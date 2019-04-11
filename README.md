# mrt-build-info

Tools for gathering Merritt build information.

## Invocation

Invocation is for the form

```
mrt-build-info <command> [flags] [URL]
```

where `<command>` is one of:

- [`jobs`](#jobs): list Jenkins jobs

and `URL` is the URL of the Jenkins server. (If not specified, `URL` defaults to `http://builds.cdlib.org/`.)

## Commands

### `jobs`

List Jenkins jobs.

By default, the `jobs` command simply lists all Jenkins jobs by name. The
flags below can be set to provide more information for each job, taken from
the last successful build.

| Short form | Flag             | Description                         |
| ---        | ---              | ---                                 |
| `-a`       | `--artifacts`    | list artifacts                      |
| `-b`       | `--build`        | show info for last successful build |
| `-r`       | `--repositories` | list repositories                   |

If any of these flags are set, output will be in the form of a tab-separated
table, with header.

Sample output:

```
$ mrt-build-info jobs -abrv

Job Name	Repository	Build	SHA Hash	Artifacts
cdl-zk-queue	https://github.com/CDLUC3/cdl-zk-queue.git	46	9defbdaff6220d6b3ed2368f74bddbc974d524bf	org.cdlib.mrt:cdl-zk-queue:0.2-20190305.190419-44 (jar, cdl-zk-queue-0.2-SNAPSHOT.jar)
git-core2	https://github.com/CDLUC3/mrt-core2.git	16	88f6022e3e622b5aef4fac43be3cc27eedf8f714	org.cdlib.mrt:mrt-core:2.0-SNAPSHOT (jar, mrt-core-2.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-core-util:2.0-SNAPSHOT (pom, pom.xml), org.cdlib.mrt:mrt-core-utilinit:2.0-SNAPSHOT (pom, pom.xml), org.cdlib.mrt:mrt-jena:2.0-SNAPSHOT (jar, mrt-jena-2.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-json:2.0-SNAPSHOT (jar, mrt-json-2.0-SNAPSHOT.jar)
git-dataONE	https://github.com/CDLUC3/mrt-dataONE.git	1	c803d9999e2107f9d96d74470026ad26e38eed8a	org.cdlib.mrt:mrt-dataonesrc:1.0-SNAPSHOT (jar, mrt-dataonesrc-1.0-SNAPSHOT.jar)
git-mrt-xoai	https://github.com/CDLIB/mrt-conf-prv.git	8	582807570c5c74786b87a838a95d055d2c1f8aea	com.lyncode:xoai:4.1.1-SNAPSHOT (pom, pom.xml), com.lyncode:xoai-common:4.1.1-SNAPSHOT (jar, xoai-common-4.1.1-SNAPSHOT.jar), com.lyncode:xoai-data-provider:4.1.1-SNAPSHOT (jar, xoai-data-provider-4.1.1-SNAPSHOT.jar), com.lyncode:xoai-service-provider:4.1.1-SNAPSHOT (jar, xoai-service-provider-4.1.1-SNAPSHOT.jar)
git-mrt-zoo	https://github.com/CDLUC3/mrt-zoo.git	11	275773ebcc558c2fab74d30353362b7fb5656d76	org.cdlib.mrt:mrt-zoopub-src:1.0-SNAPSHOT (jar, mrt-zoopub-src-1.0-SNAPSHOT.jar)
Merritt Development Submission (Full Stack Test)			5		
Merritt Production Submission (Full Stack Test)			5		
Merritt Stage Submission (Full Stack Test)			6		
mrt-build-audit	https://github.com/CDLUC3/mrt-audit	42	cbfd73970238357ebce37c6837bef1eb22260ff9	org.cdlib.mrt:mrt-audit:1.0-SNAPSHOT (pom, mrt-audit-1.0-SNAPSHOT.pom), org.cdlib.mrt:mrt-auditconfpub:1.0-SNAPSHOT (jar, mrt-auditconfpub-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-auditsrcpub:1.0-SNAPSHOT (jar, mrt-auditsrcpub-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-auditwarpub:1.0-SNAPSHOT (war, mrt-auditwarpub-1.0-SNAPSHOT.war)
mrt-build-inv	https://github.com/CDLUC3/mrt-inventory	39	a095512dfb186c6353b9eff1948374b077718ea9	org.cdlib.mrt:mrt-inventory:1.0-SNAPSHOT (pom, mrt-inventory-1.0-SNAPSHOT.pom), org.cdlib.mrt:mrt-inventoryconf:1.0-SNAPSHOT (jar, mrt-inventoryconf-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-inventorysrc:1.0-SNAPSHOT (jar, mrt-inventorysrc-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-invwar:1.0-SNAPSHOT (war, mrt-invwar-1.0-SNAPSHOT.war)
mrt-build-mysql	https://github.com/CDLIB/mrt-conf-prv.git	8	582807570c5c74786b87a838a95d055d2c1f8aea	org.cdlib.mrt:mrt-confmysql:dev-1.0-SNAPSHOT (jar, mrt-confmysql-dev-1.0-SNAPSHOT.jar)
mrt-build-oai	https://github.com/CDLUC3/mrt-oai	19	782e972109d7ad84da07b9c0332c474e13ad97ac	org.cdlib.mrt:mrt-oaiconfpub:1.0-SNAPSHOT (jar, mrt-oaiconfpub-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-oaigit:1.0-SNAPSHOT (pom, mrt-oaigit-1.0-SNAPSHOT.pom), org.cdlib.mrt:mrt-oaisrcpub:1.0-SNAPSHOT (jar, mrt-oaisrcpub-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-oaiwar:1.0-SNAPSHOT (war, mrt-oaiwar-1.0-SNAPSHOT.war)
mrt-build-replic	https://github.com/CDLUC3/mrt-replic	39	e3a17e4b63f1f35aa421e88159dfb488b383ff3f	org.cdlib.mrt:mrt-replicationconf:1.0-SNAPSHOT (jar, mrt-replicationconf-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-replicationsrc:1.0-SNAPSHOT (jar, mrt-replicationsrc-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-replicationwar:1.0-SNAPSHOT (war, mrt-replicationwar-1.0-SNAPSHOT.war), org.cdlib.mrt:mrt-replicgit:1.0-SNAPSHOT (pom, mrt-replicgit-1.0-SNAPSHOT.pom)
mrt-build-s3	ssh://git@github.com/cdlib/mrt-conf-prv.git	51	691770b9182d3870c85aba8ca776c0a3e85aa57e	org.cdlib.mrt:mrt-confs3:1.0-SNAPSHOT (jar, mrt-confs3-1.0-SNAPSHOT.jar)
mrt-build-store	https://github.com/CDLIB/mrt-conf-prv.git	85	691770b9182d3870c85aba8ca776c0a3e85aa57e	org.cdlib.mrt:mrt-confstore:prod-1.0-SNAPSHOT (jar, mrt-confstore-prod-1.0-SNAPSHOT.jar)
mrt-build-sword	https://github.com/CDLUC3/mrt-sword	10	1a01c30995d263057c4c9f486527cd25edbb28b5	org.cdlib.mrt:mrt-swordconfpub:1.0-SNAPSHOT (jar, mrt-swordconfpub-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-swordpub:1.0-SNAPSHOT (pom, mrt-swordpub-1.0-SNAPSHOT.pom), org.cdlib.mrt:mrt-swordsrcpub:1.0-SNAPSHOT (jar, mrt-swordsrcpub-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-swordwarpub:1.0-SNAPSHOT (war, mrt-swordwarpub-1.0-SNAPSHOT.war)
mrt-cloudhost-pub	https://github.com/CDLUC3/mrt-cloudhost-pub.git	19	ad697a2acfdd3af594fa359b278045e2d789f43d	org.cdlib.mrt:mrt-cloudhost:1.0-SNAPSHOT (pom, pom.xml), org.cdlib.mrt:mrt-cloudhostconfpub:1.0-SNAPSHOT (jar, mrt-cloudhostconfpub-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-cloudhostjetty:1.0-SNAPSHOT (jar, mrt-cloudhostjetty-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-cloudhostsrcpub:1.0-SNAPSHOT (jar, mrt-cloudhostsrcpub-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-cloudhostwar:1.0-SNAPSHOT (war, mrt-cloudhostwar-1.0-SNAPSHOT.war)
mrt-ingest-dev	https://github.com/CDLUC3/mrt-ingest	535	a572b8c116ef56936ffd866749bd667d29a36fdd	org.cdlib.mrt:mrt-ingest:1.0-SNAPSHOT (pom, pom.xml), org.cdlib.mrt:mrt-ingestconf:1.0-SNAPSHOT (jar, mrt-ingestconf-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-ingestinit:1.0-SNAPSHOT (pom, pom.xml), org.cdlib.mrt:mrt-ingestsrc:1.0-SNAPSHOT (jar, mrt-ingestsrc-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-ingestwar:1.0-SNAPSHOT (war, mrt-ingestwar-1.0-SNAPSHOT.war)
mrt-ingest-stage	https://github.com/CDLUC3/mrt-ingest	163	a572b8c116ef56936ffd866749bd667d29a36fdd	org.cdlib.mrt:mrt-ingest:1.0-SNAPSHOT (pom, pom.xml), org.cdlib.mrt:mrt-ingestconf:1.0-SNAPSHOT (jar, mrt-ingestconf-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-ingestinit:1.0-SNAPSHOT (pom, pom.xml), org.cdlib.mrt:mrt-ingestsrc:1.0-SNAPSHOT (jar, mrt-ingestsrc-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-ingestwar:1.0-SNAPSHOT (war, mrt-ingestwar-1.0-SNAPSHOT.war)
mrt-jetty-cloudhost	https://github.com/CDLUC3/mrt-cloudhost-pub.git	9	ad697a2acfdd3af594fa359b278045e2d789f43d	org.cdlib.mrt:mrt-cloudhost:1.0-SNAPSHOT (pom, pom.xml), org.cdlib.mrt:mrt-cloudhostjetty:1.0-SNAPSHOT (jar, mrt-cloudhostjetty-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-cloudhostsrcpub:1.0-SNAPSHOT (jar, mrt-cloudhostsrcpub-1.0-SNAPSHOT.jar)
mrt-s3-pub	https://github.com/CDLUC3/mrt-s3-pub.git	81	62739213612e3eeb59cfb408dc88968d8ba95a81	org.cdlib.mrt:mrt-nodetest:1.0-SNAPSHOT (jar, mrt-nodetest-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-openstackpub:1.0-SNAPSHOT (jar, mrt-openstackpub-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-s3:1.0-SNAPSHOT (pom, pom.xml), org.cdlib.mrt:mrt-s3srcpub:1.0-SNAPSHOT (jar, mrt-s3srcpub-1.0-SNAPSHOT.jar)
mrt-store-pub	https://github.com/CDLUC3/mrt-store.git	93	af174ac555758a1c639a7a3da39e022d9fdbf3a6	org.cdlib.mrt:mrt-storepub:1.0-SNAPSHOT (pom, mrt-storepub-1.0-SNAPSHOT.pom), org.cdlib.mrt:mrt-storepub-src:1.0-SNAPSHOT (jar, mrt-storepub-src-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-storewar:1.0-SNAPSHOT (war, mrt-storewar-1.0-SNAPSHOT.war)
mrt-test	ssh://git@github.com/cdlib/mrt-test.git	3	f0fe10c8c247d6c07a361489f163f64623e48a21	org.cdlib.mrt:mrt-s3test:1.0-SNAPSHOT (jar, mrt-s3test-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-test:1.0-SNAPSHOT (pom, pom.xml)
test-gittest	https://github.com/dloy/mrt-gittest.git	4	021c01abecdd3add2443a53e4487635be8554ff9	org.cdlib.mrt:mrt-gittest:1.0-SNAPSHOT (pom, pom.xml), org.cdlib.mrt:mrt-gittestconf:stage-1.0-SNAPSHOT (jar, mrt-gittestconf-stage-1.0-SNAPSHOT.jar), org.cdlib.mrt:mrt-gittestemb:1.0-SNAPSHOT (jar, mrt-gittestemb-1.0-SNAPSHOT.jar)
```


#### Additional flags

The `jobs` command supports the following additional flags:

| Short form | Flag           | Description                    |
| ---        | ---            | ---                            |
| `-h`       | `--help`       | help for jobs                  |
| `-v`       | `--verbose`    | verbose output                 |
| `-l`       | `--log-errors` | log non-fatal errors to stderr |


