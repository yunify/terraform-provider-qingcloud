# Proposal: Transfer Ownership 

目前 [yunify/qingcloud-terraform-provider](https://github.com/yunify/qingcloud-terraform-provider) 是从 [CuriosityChina/terraform-qingcloud](https://github.com/CuriosityChina/terraform-qingcloud) fork 的。

我和 @magicshui 商量了希望可以把  [CuriosityChina/terraform-qingcloud](https://github.com/CuriosityChina/terraform-qingcloud) 的 ownership transfer 到 [yunify/qingcloud-terraform-provider](https://github.com/yunify/qingcloud-terraform-provider) 下面，而不是以 fork 的方式。

这样做有几个好处：

+ 方便项目 issue track 和 pull request，社区更友好。
+ Qingcloud 有维护者，项目维护会更活跃，bug fix 和 添加 new feature 会更及时。
+ CuriosityChina 这边使用出问题的时候，可以像使用 qingcloud-sdk-go 或者 qingstor-sdk-go 一样提 issue 或者 pull request 就行。
+ Qingcloud 才刚开始维护这个项目，我们可以尽早把这件事情办了，后期省很多麻烦。