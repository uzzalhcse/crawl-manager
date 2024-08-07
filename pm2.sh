#! /bin/bash
git pull
cd /root/crawl-manager/backend
go build -o crawl-manager-backend
pm2 restart crawl-manager-backend

cd /root/crawl-manager/frontend
npm install
npm run build
pm2 restart crawl-manager-frontend