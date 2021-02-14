<p align="center">
    <img src="assets/logo.svg" width="275px">
</p>

<h1 align="center">MailWhale</h1>
<h3 align="center">A <i>bring-your-own-SMTP-server</i> mail relay</h3>

<p align="center">
<img src="https://badges.fw-web.space/github/license/muety/mailwhale">
<a href="https://saythanks.io/to/n1try" target="_blank"><img src="https://badges.fw-web.space/badge/SayThanks.io-%E2%98%BC-1EAEDB.svg"></a>
<a href="https://wakapi.dev" target="_blank"><img src="https://badges.fw-web.space/endpoint?url=https://wakapi.dev/api/compat/shields/v1/n1try/interval:any/project:mailwhale&color=blue"></a>
</p>

## 📄 Description

Being a web developer, chances are high that at some point you need to teach your application how to send mails.
Essentially, there are two options. Either you use a **professional mail sending service**
like [Mailgun](https://mailgun.com), [SendGrid](https://sendgrid.com), [SMTPeter](https://smtpeter.com) and the like or
you **include an SMTP client library** to your software and **plug your own mail server**.

However, if you want the best of both worlds – that is, send mails via simple HTTP calls and with no extra complexity, 
but still use your own infrastructure – you may want to go with ✉️🐳.

You get a simple **REST API**, which you can call to send out e-mail. You can plug your self-hosted SMTP server, as well as Google Mail or **literally any other e-mail provider**.

Stay tuned, there is a lot more to come.

## 🚧 Project State
The project is in a very early stage and breaking changes are likely to happen. We'd recommend to not yet use this in production or at least expect non-trivial effort required to upgrade to a new version.

## 📦 Installation
```bash
# 1. Clone repo
$ git clone https://github.com/muety/mailwhale.git

# 2. Adapt config to your needs, i.e. set your SMTP server and credentials, etc.
$ cp config.default.yml config.yml
$ vi config.yml

# 3. Run it
$ GO111MODULE=on go build
$ ./mailwhale
```

## ⌨️ Usage

MailWhale has the notion of **clients**, which are applications allowed to access the API to send mails or manage other
clients. Once you start MailWhale for the first time, a default client is created and its credentials are printed to the
console. Remember those, as you will need them to use the API.

### Authentication
Authenticating against the API currently happens through sending HTTP basic auth, i.e. the `Authorization` is set to the Base64-encoded version of `<CLIENT_NAME>:<PASSWORD>`, for instance:
```
Authorization: Basic cm9vdDpmOTE0MjM3OS0wZDVlLTQ2OTItYWJiNy1kNDJmYmEyYTJmZWYK  
```

### Create new client application

```bash
$ curl -XPOST \
  -u 'root:<your_api_key>' \
  -H 'content-type: application/json' \
  --data '{
    "name": "my-cool-app",
    "permissions": [ "send_mail" ],
    "default_sender": "John Doe <john.doe@example.org>"
  }' \
  http://localhost:3000/api/client

# Response (201 Created):
# {
#     "name": "my-cool-app",
#     "permissions": ["send_mail"],
#     "default_sender": "John Doe <john.doe@example.org>",
#     "allowed_senders": ["john.doe@example.org"],
#     "api_key": "411f9f30-3dfd-4b94-b427-2345e6e84677"
# }
```

### Send an HTML mail (synchronously)

```bash
$ curl -XPOST \
  -u 'root:<your_api_key>' \
  -H 'content-type: application/json' \
  --data '{
      "from": "John Doe <john.doe@example.org>",
      "to": ["Jane Doe <jane@doe.com>"],
      "subject": "Dinner tonight?",
      "html": "<h1>Hey you!</h1><p>Wanna have dinner tonight?</p>"
  }' \
  http://localhost:3000/api/mail
```

You can also a `text` field instead, to send a plain text message.

## 🚀 Features (planned)

Right now, this app is very basic. However, there are several cool features on our roadmap.

* **Mail Templates:** Users will be able to create complex (HTML) templates or presets for their mails, which can then
  be referenced in send requests.
* **Bound handling:** Ultimately, we want to offer the ability to plug an IMAP server in addition, to get notified about
  bounced / undelivered mails.
* **Statistics:** There will be basic statistics about when which client has sent how many mails, how many were
  successful or were bounced, etc.
* **Web UI:** A nice-looking web UI will make client- and template management easier.
* **Client libraries:** To make the developer experience even smoother, client SDKs for different programming languages will we added some time.
* **Minor enhancements:** IPv6- and TLS support, API documentation, ...

## 📓 License

MIT