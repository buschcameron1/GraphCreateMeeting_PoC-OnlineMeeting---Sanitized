**WARNING: ANY USE BY YOU OF THE CODE PROVIDED IN THIS EXAMPLE IS AT YOUR OWN RISK.**
**Microsoft provides this sample code "as is" without warranty of 'any kind, either express or implied, including but not limited to the implied warranties of 'merchantability and/or fitness for a particular purpose.**

**Example code for creating an online meeting and sending an email invite with a meeting template with Graph API**

In order to set this up you will need an App Reg configured with Calendars.ReadWrite.All and OnlineMeetings.ReadWrite.All application permissions.
You will need to set up a meeting template (https://learn.microsoft.com/en-us/microsoftteams/create-custom-meeting-template) and obtain the ID from the URL (including the customtemplate_ portion).
You will need to create and Teams application policy (https://learn.microsoft.com/en-us/graph/cloud-communication-online-meeting-application-access-policy) and set it for global (unless targeting specific users).
