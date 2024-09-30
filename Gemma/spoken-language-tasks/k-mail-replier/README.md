# Spoken Language Tasks Assistant with Gemma 

This tutorial walks you through setting up, running, and extending a spoken 
language task application built with Gemma and Python. The application provides
a basic web user interface that you can modify to fit your needs. The application
is built to generate replies to customer emails for a fictitious Korean bakery,
and all the language input and output is handled entirely in Korean. You can use
this application pattern with any language and any business task that uses text
input and text output.

## Project setup

These instructions walk you through getting this project set up for
development and testing. The general steps are installing some prerequisite
software, cloning the project from the code repository, setting a few environment 
variables, and running the configuration installation.

### Install the prerequisites

This project uses Python 3 and Python Poetry to manage packages and
run the application. The following installation instructions are for a Linux
host machine.

To install the required software:

*  Install Python 3 and the `venv` virtual environment package for Python.
<pre>
sudo apt update
sudo apt install git pip python3-venv
</pre>

### Clone and configure the project

Download the project code and use the Poetry installation command to download
the required dependencies and configure the project. You need
[git](https://git-scm.com/) source control software to retrieve the
project source code.

To download the project code:

1.  Clone the git repository using the following command.
<pre>
git clone https://github.com/google-gemini/gemma-cookbook.git
</pre>
1.  Optionally, configure your local git repository to use sparse checkout,
    so you have only the files for the project.
<pre>
cd gemma-cookbook/
git sparse-checkout set Gemma/spoken-language-tasks/
git sparse-checkout init --cone
</pre>

To install the Python libraries:

1.  Configure and activate Python virtual environment (venv) for this project:
<pre>
python3 -m venv venv
source venv/bin/activate
</pre>
1.  Install the required Python libraries for this project using the {{setup_python}} script.
<pre>
./setup_python.sh
</pre>

### Set environment variables

Set a few environment variables that are required to allow this code
project to run, including a Kaggle user name and Kaggle token key.
You must have a Kaggle account and request access to the Gemma model.

You add your Kaggle Username and Kaggle Token Key to two `.env` files, 
which are read by the web application and the tuning program, respectively.

Caution: Treat your Kaggle Token Key like a password and protect it appropriately.
Don't embed your key in publicly published code.

To set the environment variables:

1.  Obtain your Kaggle username and your token key by following the instructions
    in the [Kaggle documentation](https://www.kaggle.com/docs/api#authentication)
1.  Get access to the Gemma model by following  the *Get access to Gemma* 
    instructions in the [Gemma Setup](/gemma/docs/setup#get-access) page.
1.  Create environment variable files for the project, by creating a
    `.env` text file at *each* these location in your clone of the project:
<pre>
k-mail-replier/k_mail_replier/.env
k-gemma-it/.env
</pre>
1.  After creating the `.env` text files, add the following settings to **both** files:
<pre>
KAGGLE_USERNAME=&lt;YOUR_KAGGLE_USERNAME_HERE&gt;
KAGGLE_KEY=&lt;YOUR_KAGGLE_KEY_HERE&gt;
</pre>

### Run and test the application

1.  In a terminal window, navigate to the `spoken-language-tasks/k-mail-replier/k_mail_replier/`
    directory.
<pre>
cd spoken-language-tasks/k-mail-replier/
</pre>
1.  Run the application using the `run_flask_app.sh` script:
<pre>
./run_flask_app.sh
</pre>
