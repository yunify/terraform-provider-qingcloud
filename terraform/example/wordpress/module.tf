module "qingcloud" {
  source = "./modules/qingcloud"
}

module "wordpress" {
  source                        = "./modules/wordpress"
  wordpress_instance_public_ip  = "${module.qingcloud.public_ip}"
  wordpress_public_ssh_port     = "${module.qingcloud.wordpress_instance_ssh_port}"
  wordpress_instance_private_ip = "${module.qingcloud.wordpress_instance_private_ip}"

  mysql_instance_public_ip  = "${module.qingcloud.public_ip}"
  mysql_public_ssh_port     = "${module.qingcloud.mysql_instance_ssh_port}"
  mysql_instance_private_ip = "${module.qingcloud.mysql_instance_private_ip}"
  mysql_password            = "${var.wordpress_db_password}"
}
