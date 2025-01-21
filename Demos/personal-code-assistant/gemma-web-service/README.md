# gemma-web-service

A simple implementation of a Gemma web service using Python, Keras, JAX and the 
FastAPI library.

### Install required software

This project uses Python 3 and Virtual Environments (`venv`) to manage packages
and run the application. The following installation instructions are for a Linux
host machine.

To install the required software:

*  Install Python 3 and the `venv` virtual environment package for Python:

        sudo apt update
        sudo apt install git pip python3-venv

#### Install Python libraries

Install the Python libraries with the `venv` Python virtual environment
activated to manage Python packages and dependencies. Make sure you activate the
Python virtual environment *before* installing Python libraries with the `pip`
installer. For more information about using Python virtual environments, see the
[Python venv](https://docs.python.org/3/library/venv.html) documentation.

To install the Python libraries:

1.  In a terminal window, navigate to the `gemma-web-service` directory:

        cd Gemma/personal-code-assistant/gemma-web-service/

1.  Configure and activate Python virtual environment (venv) for this project:

        python3 -m venv venv
        source venv/bin/activate

1.  Install the required Python libraries for this project using the
    `setup_python` script:

        ./setup_python.sh

Tip: On Linux operating systems, you may need to allow execution of the bash
script by running the command `chmod +x setup_python.sh`.

#### Set environment variables

This project requires a few environmental environment variables to run,
including a Kaggle username and a Kaggle API token. You must have a Kaggle
account and request access to the Gemma models to be able to download them. For
this project, you add your Kaggle Username and Kaggle API token to two `.env`
files, which are read by the web application and the tuning program,
respectively.

Caution: Treat your Kaggle API token like a password and protect it
appropriately. Don't embed your key in publicly published code.

To set the environment variables:

1.  Obtain your Kaggle username and your token key by following the instructions
    in the [Kaggle documentation](https://www.kaggle.com/docs/api#authentication).
1.  Get access to the Gemma model by following the *Get access to Gemma*
    instructions in the [Gemma Setup](/gemma/docs/setup#get-access) page.
1.  Create an environment variable file for the project, by creating a
    `.env` text file at this location in your clone of the project:
<pre>
personal-code-assistant/gemma-web-service/.env
</pre>
1.  After creating the `.env` text file, add the following settings to it:

        KAGGLE_USERNAME=<YOUR_KAGGLE_USERNAME_HERE>
        KAGGLE_KEY=<YOUR_KAGGLE_KEY_HERE>

### Run and test the application

Once you have completed the installation and configuration of the project, run
the web application to confirm that you have configured it correctly. You should
do this as a baseline check before editing the project for your own use.

To run and test the project:

1.  In a terminal window, navigate to the `gemma-web-service` directory:

        cd personal-code-assistant/gemma-web-service/

1.  Run the application using the `run_service` script:

        ./run_service.sh

1.  After starting the web service, the program code lists a URL where
    you can access the service. Typically, this address is:

        http://localhost:8000/

1.  Test the service by running the `test_post` script:

        ./test/test_post.sh