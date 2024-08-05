# Baseline Manager

A sample project to import and export the yaml configuration file of https://github.com/s3studio/cloud-bench-checker, and display important properties in the browser.

[![Node.js](https://github.com/S3Studio/cloud-bench-checker/actions/workflows/baseline_manager_nodejs_test.yml/badge.svg)](https://github.com/S3Studio/cloud-bench-checker/actions/workflows/baseline_manager_nodejs_test.yml)

## Feature
* All data is stored in localStorage of the browser and kept local, with no remote transmission

## DISCLAIMER
This project is a sample project, so related issues would have a low priority.

## Quick Start
1. Clone this repo.
1. Install npm packages: `npm install`
1. Run as a development preview: `npm run preview`
1. Access port 3000 from a browser

## Known issues
* Changing the theme of highlight.js is not yet supported, so it is not very nice to view the raw yaml configuration under the dark theme.
* https://github.com/S3Studio/cloud-bench-checker/issues/2

## Thanks to
[Vue 3](https://v3.vuejs.org/), [Vuetify 3](https://vuetifyjs.com/en/), and all other supporting libraries
