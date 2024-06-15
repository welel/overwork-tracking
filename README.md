# Overwork Tracking App

The Overwork Tracking App is a command-line tool designed to help you monitor your daily work hours and manage your productivity effectively.

## Features

- **Record Working Hours**: Log the hours you worked each day.
- **Change Daily Goals**: Update your work hour goal for today.
- **View History**: See a history of your recorded work hours and goals.
- **Storage**: Data is stored in a JSON file.

## Usage

Once the application is running, follow the on-screen prompts to interact with the app:

- **Main Menu**: Choose options to record working hours, change daily goals, or view history.
- **Recording Hours**: Enter the hours worked in the format `HH:MM`.
- **Changing Goals**: Update the required work hours for the day in the format `HH:MM`.
- **Viewing History**: See a detailed log of recorded work hours, including overwork calculations.

## Main menu

```
-----------------------
Work Today:     08:00
Overwork:       01:15

1. Record Working Hours
2. Change Need Work
3. Print History
-----------------------
Select an option: 
```

## History Table

```
________________________________________
| Date  | Worked | Need work | Overwork |
|-------+--------+-----------+----------|
| 22.05 | 06:15  | 07:00     | -00:45   |
| 23.05 | 08:25  | 07:00     |  01:25   |
| 24.05 | 07:20  | 07:00     |  00:20   |
|       |        |           |          |
|       |        |           |          |
| 27.05 | 05:50  | 07:00     | -01:10   |
| 28.05 | 07:30  | 07:00     |  00:30   |
| 29.05 | 06:20  | 07:00     | -00:40   |
| 30.05 | 05:50  | 07:00     | -01:10   |
| 31.05 | 06:30  | 07:00     | -00:30   |
|_______|________|___________|__________|
```
