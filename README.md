# restic-plus

This is a simple wrapper around the normal [Restic](https://github.com/restic/restic) binary.

## Additions

* YAML based configuration file, see [here](restic-plus.yaml.example).
* ZFS datasets and ZFS zvols support while leveraging ZFS snapshots to garuantuee consistent backups.

## Restrictions

* Currently is opinionated and only works with SFTP.

## Usage

* `restic-plus backup`: Run backups
* `restic-plus cron`: Run backups and do cleanup afterwards
* `restic-plus -- xxx`: Forward to original restic binary
