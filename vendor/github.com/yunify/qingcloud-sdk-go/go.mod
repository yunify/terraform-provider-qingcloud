module github.com/yunify/qingcloud-sdk-go

go 1.13

require (
	github.com/DATA-DOG/godog v0.10.0
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.6.1
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/DATA-DOG/godog v0.10.0 => github.com/cucumber/godog v0.7.9

replace github.com/golang/lint v0.0.0-20201208152925-83fdc39ff7b5 => golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5

replace golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 => github.com/golang/lint v0.0.0-20201208152925-83fdc39ff7b5
