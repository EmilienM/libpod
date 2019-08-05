% podman-system-migrate(1)

## NAME
podman\-system\-migrate - Migrate container to the latest version of podman

## SYNOPSIS
** podman system migrate**

## DESCRIPTION
** podman system migrate** migrates containers to the latest podman version.

**podman system migrate** takes care of migrating existing containers to the latest version of podman if any change is necessary.

"Rootless Podman uses a pause process to keep the unprivileged
namespaces alive. This prevents any change to the `/etc/subuid` and
`/etc/subgid` files from being propagated to the rootless containers
while the pause process is running.

For these changes to be propagated, it is necessary to first stop all
running containers associated with the user and to also stop the pause
process and delete its pid file.  Instead of doing it manually, `podman
system migrate` can be used to stop both the running containers and the
pause process. The `/etc/subuid` and `/etc/subgid` files can then be
edited or changed with usermod to recreate the user namespace with the
newly configured mappings.

## SYNOPSIS
**podman system migrate**

## SEE ALSO
`podman(1)`, `libpod.conf(5)`, `usermod(8)`

## HISTORY
April 2019, Originally compiled by Giuseppe Scrivano (gscrivan at redhat dot com)
