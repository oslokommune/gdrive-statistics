# gdrive-statistics

A service for calculating usage statistics of folders in Google Drive.

To run (you need to be an gsuite admin):

* Follow step 1 here: https://developers.google.com/admin-sdk/reports/v1/quickstart/go
  * Also enable "Drive API"
  * Alternatively, you can manage this in your https://console.cloud.google.com/apis (you might need a project first)
* Download credentials.json to ~/.google-credentials.json
* Find your Google Drive Id: https://umzuzu.com/blog/2019/9/30/how-to-get-the-id-of-a-google-shared-drive-formerly-team-drives
* To run: `GOOGLE_DRIVE_ID=yourGoogleDriveId make run`
