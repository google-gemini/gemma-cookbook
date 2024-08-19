# Contributing to the Gemma Cookbook

We would love to accept your patches and contributions to the Gemma Cookbook. We are excited that you are considering donating some of your time, and this guide will help us be respectful of that time.

# Before you send anything

## Sign our contributor agreement

All contributions to this project must be accompanied by a [Contributor License Agreement](https://cla.developers.google.com/about) (CLA). You (or your employer) retain the copyright to your contribution; this simply gives us permission to use and redistribute your contributions as part of the project.

If you or your current employer have already signed the Google CLA (even if it was for a different project), you probably don't need to do it again.

Visit [https://cla.developers.google.com/](https://cla.developers.google.com/) to see your current agreements or to sign a new one.

## Style guides

Before you start writing, take a look at the [technical writing style guide](https://developers.google.com/style). You don’t need to fully digest the whole document, but do read the [highlights](https://developers.google.com/style/highlights) so you can anticipate the most common feedback.

Also check out the relevant [style guide](https://google.github.io/styleguide/) for the language you will be using. These apply strictly to raw code files (e.g. *.py, *.js), though code fragments in documentation (such as markdown files or notebooks) tend to favor readability over strict adherence.

# PR checklist 

1. Commit your finished notebook with comments and clean output after finishing the following:
   * Make sure to include the setup steps at the top (you can copy from any existing notebook), including:
        * the Colab self-link to your notebook
        * how to select GPU
        * how to set up Kaggle/HF tokens
   * Include a byline at the top of the notebook with your name, social handle and/or GitHub username
   * Run ‘nbfmt’ and 'nblint' for formatting and linting ([reference](.github/workflows/notebooks.yaml))
   * Name your notebook with words separated by underscores. For example, ‘Integrate_with_Mesop.ipynb’
2. Add the notebook name and a short description in the table of contents in README.md
3. (If applicable) remove the entry you have implemented in WISHLIST.md
4. Submit for review
5. In your PR comment, let us know if you would like your contribution to be highlighted in Google’s social handle (e.g., [Google for Developers](https://x.com/googledevs?ref_src=twsrc%5Egoogle%7Ctwcamp%5Eserp%7Ctwgr%5Eauthor) Twitter)

# Making changes

## Small fixes

Small fixes, such as typos or bug fixes, can be submitted directly via a pull request.

## Large changes or a new notebook

Before you send a PR, or even write a single line, please file an [issue](https://github.com/google-gemini/gemma-cookbook/issues). There we can discuss the request and provide guidance about how to structure any content you write.

Adding a new guide often involves lots of detailed reviews and we want to make sure that your idea is fully formed and has full support before you start writing anything. Please also check the table of contents first to avoid duplicating existing work. If you want to port an existing guide across (e.g. if you have a guide for Gemma on your own GitHub), feel free to link to it in the issue.

## Things we consider

When accepting a new guide, we want to balance a few aspects:
* Originality - e.g. Is there another guide that does the same thing?
* Pedagogy - e.g. Does this guide teach something useful? Specifically for a Gemma feature?
* Quality - e.g. Does this guide contain clear, descriptive prose? Is the code easy to understand? Is there any error?
* Practicality - e.g., Is the technique used in the guide practical in the real world?

It is not crucial for a submission to be strong along all of these dimensions, but the stronger the better. Old submissions may be replaced in favor of newer submissions that exceed these properties.
