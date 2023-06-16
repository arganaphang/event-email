# Email Sender

Email Sender prototyping using Event-Driven

## How to run

- Run `docker-compose.yaml` using

```sh
docker compose up -d --build
```

- The rest of command see [Taskfile.yaml](./Taskfile.yaml)

## Todo

- [ ] (Important) Improve Consumer Using Worker Pool
- [ ] Extract config using `.env`
- [ ] Add Dockerfile into each services
  - [ ] Payment
  - [ ] Sender
- [ ] Improve Message Data
- [ ] (Research) Schema Registry
- [ ] (Dream) Build something real with this code?
