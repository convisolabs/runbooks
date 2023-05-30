# Import Vulnerabilities to Kanbanize

This tool was created to help customers import vulnerabilities from Conviso Platform into Kanbanize.

## Pre-configuration

The first thing you want to do is edit the `statuses` variable inside the script. For example, in a test Kanbanize board, we defined the columns as the status in Conviso Platform. That means that when the script will create the vulberability card inside the Kanbanize Board, the card will already be created into the right column.
For this you need to know the id's of the columns of your board. You can check this by pressing F12 in your keyboard and inspecting the elements of the columns, as you can see in the screenshot below

![image](https://github.com/convisolabs/runbooks/assets/100381905/b6ac5857-77ef-42e3-a579-14e62fc2bc7b)

So you will want to configure the following inside the script, so the cards are created in the correct column. If you want to create them all into the same column, you can simply put the same value into all fields.
For example, in this case I want the vulnerabilities with status `identified` in Conviso Platform to be created in the column **Identified** in my Kanbanize board, so I will configure it with the id `13`, as below:

```python
statuses = {
	'undefined': 12,
	'identified': 13,
	'in_progress': 14,
	'fix_accepted': 15,
	'fix_refused': 23,
	'waiting_validation': 22
}
```

## Usage

The usage of the script is quite simple. You will need the following:

| Argument | Description | 
|---|---|
| -k,--conviso_api_key=  | Conviso Platform API key |
| -p,--company_id=  | Company Id (from Conviso Platform) |
| -n,--kanbanize_api_key=  | Kanbanize API key |
| -b,--board_id=  | Kanbanize Board Id |
| -l,--lane_id=  | Kanbanize Lane Id |
| -u,--url=  | Kanbanize URL (eg: https://conviso.kanbanize.com) |
| -f,--update | Flag for updating the columns in Kanbanize (optional) |

Then, you will run the script as:

`./kanbanize_vuln_sync.py -k <conviso_api_key> -p <company_id> -n <kanbanize_api_key> -b <board_id> -l <lane_id> -u <url> [-f]`

If you pass the `-f` or `--update` flag, the script will check if the cards already exists in the board and move them to their corresponding column (according to the status in Conviso Platform).

If you don't know the `board_id` or `lane_id` in your Kanbanize Board, you can check them by pressing F12 as well (in this case, the `board_id` is `2` and `lane_id` is `3`

![image](https://github.com/convisolabs/runbooks/assets/100381905/7960faa7-8b36-43b6-835b-de264a51e0b2)
