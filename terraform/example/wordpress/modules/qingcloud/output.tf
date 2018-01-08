output "public_ip" {
  value = "${ qingcloud_vpc.foo.public_ip }"
}

output "mysql_instance_private_ip" {
  value = "${ qingcloud_instance.mysql.private_ip }"
}

output "mysql_instance_ssh_port" {
  value = "${ qingcloud_vpc_static.ssh-mysql.val1 }"
}

output "wordpress_instance_private_ip" {
  value = "${ qingcloud_instance.wordpress.private_ip }"
}

output "wordpress_instance_ssh_port" {
  value = "${ qingcloud_vpc_static.ssh-wordpress.val1 }"
}
