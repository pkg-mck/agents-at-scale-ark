#!/bin/bash
mkdir -p certs
kubectl get secret webhook-server-cert -o jsonpath='{.data.tls\.crt}' | base64 --decode > certs/tls.crt
kubectl get secret webhook-server-cert -o jsonpath='{.data.tls\.key}' | base64 --decode > certs/tls.key
kubectl get secret webhook-server-cert -o jsonpath='{.data.ca\.crt}' | base64 --decode > certs/ca.crt