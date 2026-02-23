# tandem

![](https://api.travis-ci.com/TeamDei/tandem.svg?branch=main) ![](https://img.shields.io/github/issues/teamdei/tandem) ![](https://img.shields.io/github/downloads/teamdei/tandem/latest/total)

tandem (Latin for "at last") is a **fast, cross-platform and lightweight Latin reader** to aid in the translation/reading of Latin.

<img src="https://github.com/TeamDei/tandem/raw/main/screenshots/1.png">
<img src="https://github.com/TeamDei/tandem/raw/main/screenshots/2.png">
<img src="https://github.com/TeamDei/tandem/raw/main/screenshots/3.png">

## Features

 - Blazing fast, 100% cross-platform (Mac, Windows, Linux) and lightweight terminal application
 - Integrates with Tuft University's Perseus API to provide instant morphological analyses of Latin words
 - Integrates with the Wiktionary REST API to provide instant English dictionary definitions

## Installation

Head over to our [Releases page](https://github.com/TeamDei/tandem/releases) and download a fresh release for Windows, Mac or Linux. If you prefer, and have Go installed, simply run:

```
go install github.com/teamdei/tandem
```

## Usage

tandem is a command-line application. It allows you to read Latin texts. Navigate with the arrow keys, and look up a word's analyses by hitting the enter key. This guide assumes general familiarity with command-line applications.

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
 - **ENTER:** Get morphological analyses on a Latin word (via Perseus).
 - **TAB:** Get English dictionary definitions for a Latin word (via Wiktionary).
 - **ESC:** Exit tandem or close an open panel.

#### Tuft University / Perseus API Integration

The API integration allows tandem to instantly get scholarly analyses of a Latin word and display it to you in-editor. For example, looking up the word `terra` yields the following analyses:

```
(noun) voc. fem. sg. of terra
(noun) nom. fem. sg. of terra
(noun) abl. fem. sg. of terra
```

Following are examples which showcase the full range of the API's features:

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

#### Wiktionary API Integration

The dictionary integration allows tandem to fetch English definitions for a Latin word directly from Wiktionary. It acts natively — searching for conjugated verbs or declined nouns will automatically resolve to their base dictionary lemma. For example, hitting `Tab` on the word `principium` yields:

```
===Noun===
1. a beginning, an origin, a commencement
2. a groundwork, a foundation, a principle
3. the elements, the first principles
4. the front ranks, camp headquarters
```

## Contributing

Contributions are welcomed. Make sure to test, format and lint your code with `make test && make format && make lint` before making a pull request.

## License

This software is licensed under the GNU General Public License v3.0. See `LICENSE` for more information.
