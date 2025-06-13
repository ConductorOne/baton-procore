While developing the connector, please fill out this form. This information is needed to write docs and to help other users set up the connector.

## Connector capabilities

1. What resources does the connector sync?
- Companies
- Projects
- Users

2. Can the connector provision any resources? If so, which ones? 
- Project membership

## Connector credentials 

1. What credentials or information are needed to set up the connector? (For example, API key, client ID and secret, domain, etc.)
The connector requires a Procore App's `CLIENT_ID` and `CLIENT_SECRET`, here are the docs:
https://developers.procore.com/documentation/building-apps-create-new



2. For each item in the list above: 

   * How does a user create or look up that credential or info? Please include links to (non-gated) documentation, screenshots (of the UI or of gated docs), or a video of the process. 
    1. go to https://developers.procore.com/ and create an app
    2. after creating the app, go to the app's configuration builder section, select `Data Connector Components`, and check the service account checkbox
    3. copy the app's client id and secret id
    4. enable project directory in the projects you want to provision, you can do that by going into the project's admin section and then tool settings

   * Does the credential need any specific scopes or permissions? If so, list them here. 
   - manage companies
   - manage projects
   - manage users

   * If applicable: Is the list of scopes or permissions different to sync (read) versus provision (read-write)? If so, list the difference here. 

   * What level of access or permissions does the user need in order to create the credentials? (For example, must be a super administrator, must have access to the admin console, etc.)  
   - companies management
   - project management
   - users management
