<a id="readme-top"></a>

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/stevezaluk/protoc-go-inject-tag">
    <img src="docs/images/protobuf-logo.png" alt="Logo" width="80" height="80">
  </a>

<h3 align="center">protoc-go-inject-tag</h3>

  <p align="center">
    A command line tool for injecting custom struct tags into Protobuf's
    <br />
    <a href="https://github.com/stevezaluk/protoc-go-inject-tag"><strong>Explore the docs »</strong></a>
    <br />
    <br />
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li><a href="#disclaimers">Disclaimers</a></li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#building">Building</a></li>
        <li><a href="#usage">Usage</a></li>
      </ul>
    </li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

protoc-go-inject-tag is a command line tool for injecting custom struct tags into go-lang generated protobuf messages. I originally found a need for this when I was developing MTGJSON-API
and needed a way to insert `bson` tags into the structs that `protoc` generated. This feature is not supported out of the box using `protoc` hence the need for this project to exist. Thsi was originally
created by <a href="https://github.com/favadi/protoc-go-inject-tag">Favadi</a> and I encourage you to check out his original repo

What has changed?
* Fully featured Cobra CLI with Viper config management
* Recursively parse all files in a directory without needing to use patterns
* Proper verbosity and logging with `log/slog`
* The ability to change the file extension to search for with `--file-ext`
* Set your own comment prefix using `--comment-prefix`
* Removed deprecation warning for @inject_tag

What will be included in the future?
* A refactored API to simplify the injection process
* Proper abstractions for TextArea and TagItems
* The ability to set default tags for all fields (i.e. bson=camelCase)

## Disclaimers

protoc-go-inject-tag originally created by Favadi, improved by Steven Zaluk

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

Realistically, nothing is needed to get started with this repo. All you need to do is pull the binary from the release's page and execute it to get started. However
if you do not have any protobuf messages generated then you are going to need to install both `protoc` and the go lang support `protoc-gen-go` 

You can see instructions for installing these in the "Prerequisites" section below this one

### Prerequisites

To start install `protoc` using the following commands depending on your operating system

* MacOS
  ``brew update && brew install protobuf``
* Debian Based Linux: ``apt update && apt install -y protobuf-compiler``

You can validate that this is installed using the following command: `protoc --version`

Next install `protoc-gen-go`: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

### Building

1. Clone the repo
   ```sh
   git clone https://github.com/stevezaluk/mtgjson-api.git
   ```
2. Install dependencies
   ```sh
   go get .
   ```
3. Build the project
   ```sh
    go build
   ```
4. View the help page
    ```sh
    ./protoc-go-inject-tag --help
    ```

### Usage

```
Usage:
protoc-go-inject-tag inject [flags]

Flags:
--comment-prefix string   The prefix of the comment that protoc-go-inject-tag should search for when looking for tags to inject. A @ will be prefixed to this value (default "gotags")  -f, --file-ext string         The file extensions that should be considered for injection (default ".pb.go")
-h, --help                    help for inject
-i, --input string            The input path you want to search for protobufs with
--remove-comments         Remove comments from generated protobufs
--tags strings            Additional tags that should be applied to all fields (not implemented yet)
-v, --verbose                 Enable extended verbosity in logging

Global Flags:
--config string   config file (default is $HOME/config/protoc-go-inject-tag/config.json)
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->
## License

Distributed under Apache License 2.0. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->
## Contact

Steven A. Zaluk - [@steve_zaluk](https://x.com/stevezaluk)

Project Link: [https://github.com/stevezaluk/protoc-go-inject-tag](https://github.com/stevezaluk/protoc-go-inject-tag)\n

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/stevezaluk/protoc-go-inject-tag.svg?style=for-the-badge
[contributors-url]: https://github.com/stevezaluk/protoc-go-inject-tag/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/stevezaluk/protoc-go-inject-tag.svg?style=for-the-badge
[forks-url]: https://github.com/stevezaluk/protoc-go-inject-tag/network/members
[stars-shield]: https://img.shields.io/github/stars/stevezaluk/protoc-go-inject-tag.svg?style=for-the-badge
[stars-url]: https://github.com/stevezaluk/protoc-go-inject-tag/stargazers
[issues-shield]: https://img.shields.io/github/issues/stevezaluk/protoc-go-inject-tag.svg?style=for-the-badge
[issues-url]: https://github.com/stevezaluk/protoc-go-inject-tag/issues
[license-shield]: https://img.shields.io/github/license/stevezaluk/protoc-go-inject-tag.svg?style=for-the-badge
[license-url]: https://github.com/stevezaluk/protoc-go-inject-tag/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://www.linkedin.com/in/stevezaluk/
[go-sdk-version]: https://img.shields.io/github/go-mod/go-version/stevezaluk/protoc-go-inject-tag