# Takehome Project: Appointment Scheduling

## Run

```sh
make run
```

## Test

```sh
make test
```

## Development Environment

Optionally set up pre-commit hooks and depdendencies for a faster feedback loop. Otherwise you will have to wait for CI failures.

```sh
make setup-dev-env
```

---

### Provided Instructions

- The client should be able to pick from a list of available times, and appointments for a coach; times must not overlap.
- All appointments are 30 minutes long, and should be scheduled at `:00`, `:30` minutes after the hour during business hours.
- Business hours are M-F 8am-5pm Pacific Time
- Your goal is to create an HTTP JSON API written in Go.
- You can store appointments in a file, a database or any back end storage you prefer.

### Must Have Features

- Golang programming language
- Scheduler must be accessible via REST API
- trainer list their appointments
- client list availability
- client create appointment
- business hours are M-F 8am-5pm Pacific Time
- appointments must not overlap
- appointments are 30 minutes in duration
- appointments may be scheduled at `:00` and `:30` minutes after the hour during business hours
- load static appointments.json

### Concerns

Main issue: Time. Date logic is a complex problem. It is also possible to have timezones which further complicate date comparisons. It's usually a good idea to handle dates differently depending on which level of the stack you're in. Applications can use timezones, whereas not all storage mechanisms support timezones.

- prevent resource consumption (large start to end time ranges)
- pagination of API data
- handle Daylight Savings Time
- there will be a lot of edge cases on the logic to ensure "business is open"
- API layer handles time with timezone (especially important for scheudling and Daylight Savings Time)
- Persistence layer handles time in UTC

### Out of Scope

- Customizing trainer availability (extra after hours or blocking holidays)
- Consistent error types and friendly user messages
- Authentication/Authorization
- OpenAPI (specification, documentation)
- Observability
- Proper API routes (/v1/trainers/:id/availability)

### Required API features

- list availability (user finds availability to schedule an appointment)
- create appointment (user schedule with a trainer)
- list appointments (trainer views their appointments)

### Proposed API

All times within the API layer will use the client's timezone.

Common Status Codes:

- `200 OK`
- `400 Bad Request` - Invalid data
- `500 Internal Server Error` - Server error, try again

### List Availability

List available appointment times for a trainer between two dates.

Endpoint: `GET /availability?trainer_id=&starts_at=&ends_at=`

Parameters:

- `trainer_id`
- `starts_at`
- `ends_at`

Response

```json
{
  "appointments": [
    {
      "starts_at": "2019-01-24T09:00:00-08:00",
      "ends_at": "2019-01-24T09:30:00-08:00",
      "user_id": 1,
      "trainer_id": 1
    },
  ]
}
```

#### Create Appointment

Create a new appointment.

Considerations:

- Ensure start before end time
- User & trainer must be set, but in a real app these values wouldn't be directly assignable by a user
- Could be used by user and trainer, but for this exercise, it only works for user.
- Edge case: appointment exists (user could have been idle a long time before submitting and another client already reserved the time)

Endpoint: `POST /appointments`

Request Body

```json
{
  "appointment": {
    "trainer_id": 1,
    "user_id": 2,
    "starts_at": "",
    "ends_at": "",
  }
}
```

Response:

```json
{
 "appointment": {
  "id": 1
 }
}
```

HTTP Status Codes:

- `200` OK
- `409` Conflict - time unavailable

### List Appointments

Get a list of scheduled appointments for a trainer

Considerations:

- This could be designed to work for both the user and trainer. For this exercise it is assumed to only be used by user.

Endpoint: `GET /appointments`

Parameters:

- `trainer_id`

Response:

```json
{
  "appointments": [
    {
      "trainer_id": 1,
      "user_id": 2,
      "starts_at": "2019-01-25T09:00:00-08:00",
      "ends_at": "2019-01-25T09:30:00-08:00"
    }
  ]
}
```

---

### Data Model

- Database times will be saved in UTC.
- The backend is sqlite3.
- Normally I would combine User & Trainer, but I saw a test record both have the value 1, so I split the two into separate tables.

#### User

Fields:

- `id`
- `name`

#### Trainer

- `id`
- `name`

#### Appointment

Fields:

- `trainer_id` ForeignKey user.id
- `user_id` ForeignKey user.id
- `starts_at`
- `ends_at`

---

## Closing Thoughts

- I spent most of my time getting the repo layer working. I started with the test suite and built functions as I went. Actually wiring up the API was the last piece, which is why it's not polished. I did go over time to wire up the repo to all other parts of the project.
- Date & time is hard problem. Interactions between client and API, API and database can all have different techniques. There is also Daylight Savings Time to consider. Also things like what if we are scheduling across years (December 31 and Jan 5).
- I know that I missed some conversions of time along the way. If this were a real project I would build in helpers to ensure that times were converted between app layers. I would also spend time trying to find a way so that all dates were created the same way. For instance there were a lot of calls like `time.Now.In()` & `Truncate()` which can easily be missed.
- Leveraged existing personal projects for bootstrapping this API, database, and general approach.
- This project requires extensive validation, therefore I didn't use my normal pattern of using gin ShouldBind struct validation.
- GenerateAppointments limitations:
  - No support for gaps in business hours (e.g., lunch breaks, holidays).
  - No support for extended hours (partial nights/weekends).
  - TimeZone-to-UTC handling is cumbersome and should default to reliable automation (e.g., data mapper, ORM hooks, custom JSON marshaller).
  - The fixed 30-minute appointment duration is definitely going to be painful to change
- API lacks sufficient unit tests due to time constraints.
- General error handling lacks consistent error types, codes, and messages.
- Missing a robust error handling strategy to map errors between layers (gorm missing record -> 404)
- It would have been nice to auto generate documentation from OpenAPI spec.
- The Project file structure is over-engineered for the scope, but this is how I build real projects.
