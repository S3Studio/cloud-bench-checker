// Package main:
// Command line tool of cloud-bench-checker
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"
	"github.com/s3studio/cloud-bench-checker/pkg/framework"

	"github.com/cheggaaa/pb"
	"github.com/spf13/pflag"
	yaml "gopkg.in/yaml.v3"
)

var (
	confFilePath = pflag.StringP("conf-file", "c", "", "File containing configs and baselines in yaml format")
	tag          = pflag.StringSliceP("tag", "t", []string{"test"}, "Tags of which baselines to check")
	showProgress = pflag.BoolP("show-progress", "p", true, "Show progress")
)

// Add visibility management to pb.ProgressBar
type pbWrapper struct {
	bar     *pb.ProgressBar
	visible bool
}

func newPb(visible bool, total int, prefix string) *pbWrapper {
	if visible {
		bar := pb.New(total).Prefix(prefix).SetRefreshRate(time.Second)
		bar.ShowCounters = true
		return &pbWrapper{bar: bar, visible: true}
	} else {
		return &pbWrapper{bar: nil, visible: false}
	}
}

func (pb *pbWrapper) Start() {
	if pb.visible && pb.bar != nil {
		pb.bar.Start()
	}
}

func (pb *pbWrapper) Increment() {
	if pb.visible && pb.bar != nil {
		pb.bar.Increment()
	}
}

func (pb *pbWrapper) Finish() {
	if pb.visible && pb.bar != nil {
		pb.bar.Finish()
	}
}

func main() {
	pflag.Parse()

	var conf def.ConfFile

	if len(*confFilePath) == 0 {
		log.Println("Required parameter conf-file is missing")
		os.Exit(-1)
	}
	if confFile, err := os.Open(*confFilePath); err != nil {
		log.Printf("Failed to open \"%s\": %v\n", *confFilePath, err)
		os.Exit(-1)
	} else {
		defer confFile.Close()

		confFileContent, err := io.ReadAll(confFile)
		if err != nil {
			log.Printf("Failed to read \"%s\": %v\n", *confFilePath, err)
			os.Exit(-1)
		}

		if err = yaml.Unmarshal(confFileContent, &conf); err != nil {
			log.Printf("Failed to load \"%s\" as yaml: %v\n", *confFilePath, err)
			os.Exit(-1)
		}
	}

	// Inspect for duplicated id of listor
	for i := range conf.Listor {
		for j := 0; j < i; j++ {
			if conf.Listor[i].Id == conf.Listor[j].Id {
				log.Printf("Found duplicate id of listor %d, the latter one will be omitted\n", i)
			}
		}
	}

	if conf.Option.PageSize >= 10 {
		framework.SetPageSize(conf.Option.PageSize)
	}
	authProvider := auth.NewAuthFileProvider(conf.Profile)

	// Create baselines with tags specified in command parameter
	var baseline []*framework.Baseline
	for _, c := range conf.Baseline {
		bAdd := false
	TagFilterLoop:
		for _, tagBaseline := range c.Tag {
			for _, tagFilter := range *tag {
				if tagBaseline == tagFilter {
					bAdd = true
					break TagFilterLoop
				}
			}
		}
		if bAdd {
			baseline = append(baseline, framework.NewBaseline(&c, authProvider, nil))
		}
	}

	// Id of listor to be created
	bar1 := newPb(*showProgress, len(baseline), "Collect listor info")
	bar1.Start()
	var idListor []int

	for _, b := range baseline {
		for _, newid := range b.GetListorId() {
			bAdd := true
			for _, exid := range idListor {
				if newid == exid {
					bAdd = false
					break
				}
			}
			if bAdd {
				idListor = append(idListor, newid)
			}
		}

		bar1.Increment()
	}
	bar1.Finish()

	// Create listor and get raw data
	bar2 := newPb(*showProgress, len(idListor), "Get data from listor")
	bar2.Start()
	var wg2 sync.WaitGroup
	wg2.Add(len(idListor))
	mapRawData := &framework.SyncMapDataProvider{}

	for _, id := range idListor {
		go func(target *framework.SyncMapDataProvider) {
			defer func() {
				bar2.Increment()
				wg2.Done()
			}()

			bFound := false
			for _, c := range conf.Listor {
				if c.Id == id {
					listor := framework.NewListor(&c, authProvider)
					rawData, err := listor.ListData()
					if err != nil {
						log.Println(err)
					} else if len(rawData) > 0 {
						target.DataMap.Store(id, rawData)
						target.CtMap.Store(id, string(c.CloudType))
					}

					bFound = true
					break
				}
			}

			if !bFound {
				log.Printf("failed to find listor with id of %d, please check the conf file\n", id)
			}
		}(mapRawData)
	}

	wg2.Wait()
	bar2.Finish()

	for _, b := range baseline {
		b.SetDataProvider(mapRawData)
	}

	// Extract prop
	bar3 := newPb(*showProgress, len(baseline), "Extract prop from data")
	bar3.Start()
	var wg3 sync.WaitGroup
	wg3.Add(len(baseline))

	listProp := make([]framework.BaselinePropList, len(baseline))
	for i, b := range baseline {
		go func(target *framework.BaselinePropList) {
			defer func() {
				bar3.Increment()
				wg3.Done()
			}()

			*target = append(*target, b.GetProp()...)
		}(&listProp[i])
	}

	wg3.Wait()
	bar3.Finish()

	// Validate prop
	bar4 := newPb(*showProgress, len(baseline), "Validate prop")
	bar4.Start()

	type overallResult struct {
		baseline *framework.Baseline
		res      []*framework.ValidateResult
	}
	result := make([]overallResult, 0, len(baseline))
	for i, b := range baseline {
		resBaseline, err := b.Validate(listProp[i])
		if err != nil {
			log.Println(err)
		} else if len(resBaseline) > 0 {
			result = append(result, overallResult{
				baseline: b, res: resBaseline,
			})
		}

		bar4.Increment()
	}
	bar4.Finish()

	// Output result
	var outputData []map[string]string
	for _, eachResBaseline := range result {
		for _, eachRes := range eachResBaseline.res {
			if eachRes.InRisk || !conf.Option.OutputRiskOnly {
				singleOutputData := map[string]string{
					"Cloud Type":    string(eachRes.CloudType),
					"Resource Id":   eachRes.Id,
					"Resource Name": eachRes.Name,
					"Actual Value":  eachRes.Value,
				}
				if eachRes.InRisk {
					singleOutputData["Resource in risk"] = "True"
				} else {
					singleOutputData["Resource in risk"] = "False"
				}

				for _, key := range conf.Option.OutputMetadata {
					value := (*eachResBaseline.baseline.GetMetadata())[key]
					singleOutputData[key] = value
				}

				outputData = append(outputData, singleOutputData)
			}
		}
	}

	if (conf.Option.OutputFormat == def.OUTPUT_FORMAT_CSV ||
		conf.Option.OutputFormat == def.OUTPUT_FORMAT_JSON) &&
		len(conf.Option.OutputFilename) > 0 {
		file, err := os.Create(fmt.Sprintf("%s.%s", conf.Option.OutputFilename, conf.Option.OutputFormat))
		if err != nil {
			log.Printf("failed to open output file: %v\n", err)
		} else {
			defer file.Close()

		OutputSwitch:
			switch conf.Option.OutputFormat {
			case def.OUTPUT_FORMAT_JSON:
				by, err := json.Marshal(outputData)
				if err != nil {
					log.Printf("failed to marshal result as json: %v\n", err)
					break
				}

				if _, err := file.Write(by); err != nil {
					log.Printf("failed to output result: %v\n", err)
					break
				}

				log.Printf("%d result(s) are outputed", len(outputData))
			case def.OUTPUT_FORMAT_CSV:
				file.WriteString("\xEF\xBB\xBF")                      // UTF8-BOM for Excel
				regNum := regexp.MustCompile(`^(\d*\.)?\d+(\.\d*)?$`) // Check numberic value

				header := []string{"Cloud Type", "Resource Id", "Resource Name", "Resource in risk", "Actual Value"}
				for _, key := range conf.Option.OutputMetadata {
					if regNum.MatchString(key) {
						key = fmt.Sprintf("\"%s\"", key) // Avoid item to be convert to integer
					}
					header = append(header, key)
				}

				w := csv.NewWriter(file)
				if err := w.Write(header); err != nil {
					log.Printf("failed to output to csv file: %v\n", err)
					break
				}

				for _, eachData := range outputData {
					line := make([]string, len(header))
					for i, key := range header {
						value := eachData[key]
						if regNum.MatchString(value) {
							value = fmt.Sprintf("\"%s\"", value) // Avoid item to be convert to integer
						}
						line[i] = value
					}

					if err := w.Write(line); err != nil {
						log.Printf("failed to output to csv file: %v\n", err)
						break OutputSwitch
					}
				}

				w.Flush()
				log.Printf("%d result(s) are outputed", len(outputData))
			}
		}
	} else {
		fmt.Printf("No valid output config in the conf file. %d result(s) waiting to be output.\n", len(outputData))
	}
}
