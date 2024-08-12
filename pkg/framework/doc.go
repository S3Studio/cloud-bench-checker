// Package framework:
// Overall management of the benchmarking process, including Baseline, Checker and Listor
//
// Explanation:
//
//   - Listor:
//
//     Used to retrieve a list of resources and their basic information from the cloud with connector.
//
//   - Checker:
//
//     Used to extract required properties and validate that they meet the requirements of benchmark guidelines.
//
//     1. Checker.GetProp:
//
//     Used to get property either by extracting it from existing data of a listor,
//     or by retrieving it via another API from the cloud if required.
//
//     2. Checker.Validate:
//
//     Used to validate the property against the benchmark and return the result.
//
//     3. NOTE:
//
//     It is useful to separate the GetProp and Validate steps into different functions to serve them
//     from different servers, or from local-side and remote-side,
//     but it is also acceptable to put them in sequence in your own code.
//
//   - Baseline:
//
//     Used to manage checkers and listors.
//     It is recommended that each baseline corresponds to a single benchmark recommendation.
package framework
