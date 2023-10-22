# Society Voting

*Online voting designed for student groups*

---

## TODO

### System

- [x] Database setup
- [x] Account provisioning and login
- [x] Guild data scraping
- [x] Events
  - [ ] Discord webhook event notifier
- [x] Use database transactions!
- [x] Change package namespace
- [ ] Save election results to dedicated table
- [x] Make election results prettier

### API

- [x] Allow admin access to non-user-specific sections of the normal API
- [x] Add `isRON` flag to `BallotEntry`
- [x] Vote validation code
- [x] Add error messages to all status code-only responses where applicable

#### User

- [x] Change display name
- [x] Stand/withdraw from election
- [x] List all elections
- [x] Display currently running election in /api/elections
- [x] Vote endpoint
- [x] Make only the main election list endpoint return the candidate list
- [ ] Ensure vote IDs are real ballot options

#### Admin

- [x] Create election
- [x] Delete election
- [x] Run election
  - [x] Add `Ballot` table 
  - [x] Accept extra ballot options
  - [x] Create ballot in setup endpoint
  - [x] Store active election
- [x] Stop and finalise election
- [x] Election status SSE
- [x] Delete user
- [x] Remove candidate
- [x] Make active election endpoint return the number of eligible voters

### Frontend
