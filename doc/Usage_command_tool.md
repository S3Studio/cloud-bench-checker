# Command tool

Perform security baseline checks across multiple clouds according to benchmark recommendations in command-line mode.

## Feature
* Support for multiple platforms including Windows, Linux and MacOS(built in release assets but not yet tested)

## Prepare
### Configuration file
* Prepare a baseline configuration file of your interest based on the [reference](Baseline.md)

or

* Select one of the template configuration files from the [directory](/template/)

### Cloud auth config
Prepare authorization information from the corresponding cloud, and store it in the format of the [reference](Auth.md).

*DISCLAIMER*:
**ALWAYS** use the *READONLY* cloud authorizations (ak/sk/ClusterRole/etc...) to be configured in the project,
and **NEVER** trust any rule provided by others, even if it is cloned or downloaded from this site.

## Run locally
### Download binary
Download the command tool for your OS from the [release page](https://github.com/S3Studio/cloud-bench-checker/releases).

### Check environment
If the name of the profile of a cloud in the conf file is `$ENV`,
remember to store authorization information in environment variables or `${HOME}/.kube/config`(for "k8s").

### Command-line argument
Run `./main -h`, and the instruction will look like this:
```
Usage of xxx/main:
  -c, --conf-file string   File containing configs and baselines in yaml format
  -p, --show-progress      Show progress (default true)
  -t, --tag strings        Tags of which baselines to check (default [test])
```

#### --conf-file, -c
The baseline configuration file prepared above. Required: true

#### --show-progress, -p
Display a progress bar showing the rate of each step of the check.

The effect of the progress bar varies between platforms.
On my test platform of Windows, it looks like this:
```
Collect listor info 3 / 3 [========] 100.00% 0s
Get data from listor 4 / 4 [========] 100.00% 2s 
Extract prop from data 3 / 3 [========] 100.00% 0s 
Validate prop 3 / 3 [========] 100.00% 0s
```

In some situations such as debugging, other information may be output to the stdout,
so it is designed to switch off the progress bar with this argument
to avoid mixing output from different sources.

#### --tag, -t
Specific the Baseline which matches the tag to be checked.

Tag is designed to making it easy for different customers, departments, etc.
to use parts of one conf file that interest them.
See the [reference](./Baseline.md#tag)

The tag argument accepts multiple values that are combined using the *OR* logic.
So the Baseline with any tag in the provided list is considered to match the argument.

### Output result
The output result is defined in the configuration file with the file name and format.

See the [reference](./Baseline.md#option)

## Run with Docker
### Docker image
The docker image is published [here](https://github.com/S3Studio/cloud-bench-checker/pkgs/container/cloud-bench-checker).

Alternatively, a custom image can be built with the [Dockerfile](/Dockerfile) using the following command:
```sh
docker build -t {image_name} -f ./Dockerfile .
```

### Command to start Docker container
Considering the template of the following command:
```sh
docker run --rm --env-file {env_file} -v {conf_file}:/app/config.conf -v {output_dir}:/app/output ghcr.io/s3studio/cloud-bench-checker:latest -t {tag}
```

#### --rm
As each check is a standalone process,
it is recommended to remove the container automatically after each run,
unless the information in the container is needed for debugging purposes.

#### --env-file
If the name of the profile of a cloud in the conf file is `$ENV`,
the authorization information can be passed to the container with `--env-file` argument.

Otherwise, the auth file must be mounted as a volume under the `.auth` subdirectory:
```sh
-v {conf_file}:/app/.auth/{conf_file}
```

Also, if the cloud type is "k8s" and the profile name is `$ENV`,
the auth file (kubeconfig) must be in the default location with the home directory `/home/nonroot/`:
```sh
-v {kube_config}:/home/nonroot/.kube/config
```

*NOTE:* Remember to add the current directory prefix "." to the name of the file
to prevent Docker from using a built-in volume instead of the file on the local disk.

#### -v {conf_file}:/app/config.conf
The Docker image uses "/app/config.conf" as the default baseline configuration file,
so only the local filename is required with no additional `--conf-file` argument.

If a different name of the conf file needs to be used in some situation,
pass the volume argument along with an `--entrypoint` argument to Docker:
```sh
docker {...} -v {conf_file}:/app/{new_name} --entrypoint /app/main ghcr.io/s3studio/cloud-bench-checker:latest -c {new_name} {...}
```

*NOTE:* Remember to add the current directory prefix "." to the name of the file
to prevent Docker from using a built-in volume instead of the file on the local disk.

#### -v {output_dir}:/app/output
The directory where the output result is saved.

The default output filename is "test.csv" for most conf files in the template directory.
It is hard to specific it to a file outside the container,
so the Docker image uses it as an alias to "output/output.csv" with symbolic link,
and a volume mounted to the output directory is required
to get the "output.csv" outside the container.

If the output file name or format is different from "test.csv",
a directory name is required and must be mounted in the same way.
See the [reference](./Baseline.md#option)

*NOTE:*
* Remember to add the current directory prefix "." to the name of the output directory
  to prevent Docker from using a built-in volume instead of the directory on the local disk.
* Remember to create the output directory in advance to bypass the permission limit
  if Docker needs be launched with `sudo` and creates the output directory with owner of "root".

#### -t {tag}
Specific the Baseline which matches the tag to be checked.

It can be omitted and the command tool will use `[test]` as its default value.

Tag is designed to making it easy for different customers, departments, etc.
to use parts of one conf file that interest them.
See the [reference](./Baseline.md#tag)

The tag argument accepts multiple values that are combined using the *OR* logic.
So the Baseline with any tag in the provided list is considered to match the argument.
