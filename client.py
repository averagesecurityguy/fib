#!/bin/env python3

import requests
import re

base = 'http://127.0.0.1:3000'
re_action = re.compile(r'action="/([0-9]+)/([0-9]+)">')
ses = ''
cur = 0
nex = 1

if __name__ == '__main__':
	resp = requests.get(base)
	body = resp.text
	m = re_action.search(body)

	if m is not None:
		ses = m.group(1)
		cur = int(m.group(2))

	next_url = '{0}/{1}/{2}'.format(base, ses, cur)
	while True:
		print(next_url)
		resp = requests.get(next_url, allow_redirects=False)
		if resp.status_code == 302:
			print(resp.headers['Location'])
			break

		nex = int(resp.headers['next'])
		cur = cur + nex

		next_url = '{0}/{1}/{2}'.format(base, ses, cur)
