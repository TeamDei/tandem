# tandem
tandem (Latin for "at last") is a **fast, cross-platform and lightweight editor** to aid in the translation/reading of Latin.

## (!) Notice
**tandem is a work in progress and undergoing heavy development by a volunteer team. It is not yet finished, but we plan to release a working v1.0.0 by Christmas. Stay tuned.**

## Features
 - Blazing fast, cross-platform and lightweight
 - Integrates with Tuft University's Perseus API to provide instant analyses of Latin words

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
