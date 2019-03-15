# a web terminal to ssh

## [Vuejs Admin Demo UserName:`admin` Password:`admin`](http://felix.mojotv.cn:2222)

## command: `felix sshw`
```bash
$ felix sshw -h
the demo website is http://felix.mojotv.cn:2222

Usage:
  felix sshw [flags]

Flags:
  -a, --addr string       listening addr (default ":2222")
  -x, --expire uint       token expire in * minute (default 1440)
  -h, --help              help for sshw
  -p, --password string   auth password (default "admin")
  -s, --secret string     jwt secret string
  -u, --user string       auth user (default "admin")

Global Flags:
      --verbose   verbose
```

```bash
$ felix sshw
use random string as jwt secret: F4@e~pVe6IO^4T9CL$~HJ~YH!o6ZzC~I
login user: admin
login password: admin
login expire in 60000 minutes
[GIN-debug] POST   /api/login                --> github.com/dejavuzhou/felix/ssh2ws/internal.Login (2 handlers)
[GIN-debug] GET    /api/ws/:id               --> github.com/dejavuzhou/felix/ssh2ws/internal.WsSsh (3 handlers)
[GIN-debug] GET    /api/ssh                  --> github.com/dejavuzhou/felix/ssh2ws/internal.SshAll (3 handlers)
[GIN-debug] POST   /api/ssh                  --> github.com/dejavuzhou/felix/ssh2ws/internal.SshCreate (3 handlers)
[GIN-debug] GET    /api/ssh/:id              --> github.com/dejavuzhou/felix/ssh2ws/internal.SshOne (3 handlers)
[GIN-debug] PATCH  /api/ssh                  --> github.com/dejavuzhou/felix/ssh2ws/internal.SshUpdate (3 handlers)
[GIN-debug] DELETE /api/ssh/:id              --> github.com/dejavuzhou/felix/ssh2ws/internal.SshDelete (3 handlers)
[GIN-debug] GET    /api/sftp/:id             --> github.com/dejavuzhou/felix/ssh2ws/internal.SftpLs (3 handlers)
[GIN-debug] GET    /api/sftp/:id/dl          --> github.com/dejavuzhou/felix/ssh2ws/internal.SftpDl (3 handlers)
[GIN-debug] GET    /api/sftp/:id/cat         --> github.com/dejavuzhou/felix/ssh2ws/internal.SftpCat (3 handlers)
[GIN-debug] GET    /api/sftp/:id/rm          --> github.com/dejavuzhou/felix/ssh2ws/internal.SftpRm (3 handlers)
[GIN-debug] GET    /api/sftp/:id/rename      --> github.com/dejavuzhou/felix/ssh2ws/internal.SftpRename (3 handlers)
[GIN-debug] GET    /api/sftp/:id/mkdir       --> github.com/dejavuzhou/felix/ssh2ws/internal.SftpMkdir (3 handlers)
[GIN-debug] POST   /api/sftp/:id/up          --> github.com/dejavuzhou/felix/ssh2ws/internal.SftpUp (3 handlers)
[GIN-debug] POST   /api/ginbro/gen           --> github.com/dejavuzhou/felix/ssh2ws/internal.GinbroGen (3 handlers)
[GIN-debug] POST   /api/ginbro/db            --> github.com/dejavuzhou/felix/ssh2ws/internal.GinbroDb (3 handlers)
[GIN-debug] GET    /api/ginbro/dl            --> github.com/dejavuzhou/felix/ssh2ws/internal.GinbroDownload (3 handlers)
[GIN-debug] GET    /api/term-log             --> github.com/dejavuzhou/felix/ssh2ws/internal.TermLogAll (3 handlers)
[GIN-debug] GET    /api/term-log/:id         --> github.com/dejavuzhou/felix/ssh2ws/internal.TermLogOne (3 handlers)
[GIN-debug] DELETE /api/term-log/:id         --> github.com/dejavuzhou/felix/ssh2ws/internal.TermLogDelete (3 handlers)
[GIN-debug] PATCH  /api/term-log             --> github.com/dejavuzhou/felix/ssh2ws/internal.TermLogUpdate (3 handlers)
[GIN-debug] GET    /api/user                 --> github.com/dejavuzhou/felix/ssh2ws/internal.UserAll (3 handlers)
[GIN-debug] POST   /api/user                 --> github.com/dejavuzhou/felix/ssh2ws/internal.UserCreate (3 handlers)
[GIN-debug] DELETE /api/user/:id             --> github.com/dejavuzhou/felix/ssh2ws/internal.UserDelete (3 handlers)
[GIN-debug] PATCH  /api/user                 --> github.com/dejavuzhou/felix/ssh2ws/internal.UserUpdate (3 handlers)
[GIN-debug] Listening and serving HTTP on :2222
```

![](/images/felix_sshw_001.png)

![](/images/felix_sshw_002.png)

![](/images/felix_sshw_003.png)

![](/images/felix_sshw_004.png)