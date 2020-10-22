# tandem
tandem (Latin for "at last") is a **fast, cross-platform and lightweight Latin reader** to aid in the translation/reading of Latin.

<img src="https://github.com/TeamDei/tandem/raw/main/screenshots/1.png">
<img src="https://github.com/TeamDei/tandem/raw/main/screenshots/2.png">

## (!) Notice
 - **tandem is currently in a beta state and usable. We plan to release a stable v1.0.0 by Christmas. Stay tuned.**

## Features
 - Blazing fast, 100% cross-platform (Mac, Windows, Linux) and lightweight terminal application
 - Integrates with Tuft University's Perseus API to provide instant analyses of Latin words

## Installation
Head over to our [Releases page](https://github.com/TeamDei/tandem/releases) and download a fresh release for Windows, Mac or Linux. If you prefer, and have Go installed, simply run:

```
go install github.com/teamdei/tandem
```

## Usage
tandem is a command-line application. This guide assumes general familiarity with command-line applications.

#### Mac/Linux
Open up Terminal in the folder where tandem was downloaded to and type in `./tandem -file "path/to/file.txt"` to open up the interactive reader on the specified file.

#### Windows
Using command prompt, type in `tandem.exe -file "path\to\file.txt"` to open up the interactive reader on the specified file.

## Navigation
tandem recognizes various keybinds which enable quick and easy navigation.

 - **h, left arrow:** Move left one word.
 - **l, right arrow:** Move right one word.
 - **j, down arrow:** Move down by one line.
 - **k, up arrow:** Move up by one line.
 - **g, home:** Move to the top.
 - **G, end:** Move to the bottom.
 - **Ctrl-F, page down:** Move down by one page.
 - **Ctrl-B, page up:** Move up by one page.
 - **ENTER:** Get analyses on a Latin word.
 - **ESC:** Exit tandem.

#### Tuft University / Perseus API Integration
The API integration allows tandem to instantly get scholarly analyses of a Latin word and display it to you in-editor. For example, looking up the word `terra` yields the following analyses:

```
(noun) voc. fem. sg. of terra
(noun) nom. fem. sg. of terra
(noun) abl. fem. sg. of terra
```

Following is a test which showcases the full range of the API's features:

```
===noun test (terra)===
(noun) voc. fem. sg. of terra
(noun) nom. fem. sg. of terra
(noun) abl. fem. sg. of terra

===pronoun test (ejus)===
(pron) (indeclform) gen. neut. sg. of is
(pron) (indeclform) gen. masc. sg. of is
(pron) (indeclform) gen. fem. sg. of is

===verb test (laudo)===
(verb) 1st person sg, act. ind. pres. of laudo

===adjective test (bonus)===
(adj) nom. masc. sg. of bonus

===adverb test (libere)===
(adv) liber
(adj) voc. neut. sg. of liber
(adj) voc. masc. sg. of liber
(verb) act. inf. pres. of libet
(verb) 2nd person sg, pass. ind. pres. of libet
(verb) 2nd person sg, pass. imperat. pres. of libet
(verb) 2nd person sg, pass. subj. pres. of libo
```

## Contributing
Contributions are welcomed. Make sure to format your code with `gofmt` before making a pull request.
