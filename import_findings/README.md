# Import findings into Conviso Platform
These scripts were created to facilitate the import of findings from external tools into Conviso Platform.

## Import findings from CSV
This is a script to import findings from a generic .csv file. You can check the `findings.csv` file with the expected format.

### Execution
To execute this script simply run the following command:

`./sca_findings_from_csv.py -k <conviso_api_key> -p <project_id> -f <csv_input_file>`

## Import findings from DependaBot
This is a script to import findings from the DependaBot tool from GitHub. You can import directly from GitHub (using your Authentication key), or you can import a .json file.

### Execution directly from GitHub
To execute this script so it connects to GitHub, fetches all DependaBot alerts, and insert them into Conviso Platform, you can run the following command:

`./insert_dependabot_findings.py -k <conviso_api_key> -p <project_id> -g <github_api_key> -o <github_owner> -r <github_repo>`

The GitHub owner and repo are simply obtainable from the GitHub URL. Example:
http://github.com/Company/Test_project

So the Owner is `Company` and the Repository is `Test_project`

### Execution from a .json file
If you already have the Dependabot alerts exported, you can simply run the following command:

`./insert_dependabot_findings.py -k <conviso_api_key> -p <project_id> -f <json_input_file>`

## Import findings from AWS Script
This is a script to import findings from a script output that will fetch alerts from AWS. 

### Execution
To execute this script, simply run:

`./findings_from_aws.py -k <conviso_api_key> -p <project_id> -f <csv_input_file>`