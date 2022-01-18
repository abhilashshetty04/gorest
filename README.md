# library
Rest API written in GO backed by mysql

Install helm

Add helm repo using:

helm repo add gorestapi https://abhilashshetty04.github.io/library/charts 

install Application using:

helm install lib gorest/library --set database.volume.storageClassName=your_storage_class_name

Without ingress resource APIs can be accessed using:

kubectl port-forward -n restapi svc/restapi 8080

Create database and table named library and books respectively in mysql pod using following steps to get started:

1. Login to mysql pod using
kubectl exec -it mysql_pod_name -n database bash

2. Login to mysql instance using:
mysql -u root --password=VMware@123

3. Create database
create database library;

4. Enter db:
use library;

5. Create table named books:
create table books(id varchar(100), name varchar(100), isbn varchar(100));

6.Insert initial values
insert into books values ("1", "Abhilash book", "ISBN-4");

Post this access localhost:8080/api/v1/books (given port forward is done) to get object. Post can also be used with json formatted input to add new entries to the table.
