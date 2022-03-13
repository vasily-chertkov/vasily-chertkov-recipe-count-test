#!/usr/bin/env python

import random
import argparse
import logging
import secrets
import os
import json
import sys

def main():
    logging.getLogger().setLevel(os.environ.get("LOG_LEVEL", "INFO"))

    parser = argparse.ArgumentParser()
    parser.add_argument('-vmc', dest='vmcount', required=True,
                        help="number of vm records generated", type=int)
    parser.add_argument("-fwc", dest='fwcount', required=True,
                        help="number of fw rules generated", type=int)
    parser.add_argument('outfile', nargs='?', type=argparse.FileType('w'),
        default=sys.stdout)

    args = parser.parse_args()

    AVAILABLE_TAGS = [
        "antivirus",
        "api",
        "ci",
        "corp",
        "dev",
        "django",
        "http",
        "https",
        "k8s",
        "loadbalancer",
        "nat",
        "reverse_proxy",
        "ssh",
        "storage",
        "windows-dc"
    ]

    TAGS_MIN_COUNT = 0
    TAGS_MAX_COUNT = 4

    data = {
        "vms": [],
        "fw_rules": []
    }

    vm_ids = {}
    for i in range(args.vmcount):
        id = None
        while id is None or id in vm_ids:
            id = secrets.token_hex(nbytes=random.randint(3, 5))

        vm_ids[id] = True

        vm = {
            "vm_id": "vm-{}".format(id),
            "name": "vm name for {}".format(id),
            "tags": random.sample(AVAILABLE_TAGS, k=random.randint(TAGS_MIN_COUNT, TAGS_MAX_COUNT))
        }
        data["vms"].append(vm)

    fw_ids = {}
    for i in range(args.fwcount):
        id = None
        while id is None or id in fw_ids:
            id = secrets.token_hex(nbytes=random.randint(3, 5))

        fw_ids[id] = True

        fw_rule = {
            "fw_id": "fw-{}".format(id),
            "source_tag": random.choice(AVAILABLE_TAGS),
            "dest_tag": random.choice(AVAILABLE_TAGS)
        }
        data["fw_rules"].append(fw_rule)

    args.outfile.write(json.dumps(data, indent=2))


if __name__ == '__main__':
    main()
