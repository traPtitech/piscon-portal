# インフラ構成
## 概要
ISUCONは過去問環境を構築するためのAMIがあるため, それを利用してAWS上で競技環境を整備するのが早いです. (以前はConoHaのVPSを使っており, conohaディレクトリはその名残です)
ここではAWS上のインフラ構成を説明します. 具体的なAWSの設定については [infra-setup.md](./infra-setup.md) などを見てください.

## 構成
アーキテクチャ図を示します.

![](./figure/infra.png)

### AWS Cloud
AWSはクラウドサービスの一つです. 
物理サーバーやネットワーク機器を自前で用意することなく, 柔軟にサーバー環境を構築することができます.

### Region
AWSは世界各地にデータセンターを置き, 地域単位で分割してサービスを提供しています. この分割単位を**リージョン** (Region) と呼びます.
リージョンは好きに選ぶことができますが, 遅延などの問題を避けるため, アジアパシフィック (東京) リージョンを利用することを推奨します. 

### VPC
AWS上の仮想ネットワークです.
VPC内でサブネットやルーティング,  ゲートウェイの設定を行うことで, ネットワークを構成します.

### IGW
(図の紫の橋のマーク)
インターネットゲートウェイの略です. VPCとインターネットとの通信を可能にします.

### Subnet
IPアドレスの範囲で指定される, VPC内の小さなネットワークです.
このうち, インターネットゲートウェイへのルートがあるものを**パブリックサブネット**, ないものを**プライベートサブネット**と呼びます. 定義の通り, パブリックサブネットはインターネットゲートウェイを通じて, VPCの外からアクセスすることができます.

### EC2
Elastic Compute Cloud の略で, インスタンスと呼ばれる仮想的なサーバーの構築を行うことができるサービスです.
PISCONでは

- ポータル兼ベンチマーカー用インスタンス
- 競技用インスタンス (各チーム3台)

を用意します.

## ユースケース
1. ユーザーがポータルにアクセスする.
2. ポータル上でチームを作成し, インスタンスを起動する.
3. ポータル上でベンチマークの実行をリクエストする.
4. アプリケーションがベンチマーカーを起動する. ベンチマーカーは指定された競技用インスタンスにベンチマークを行い, 結果を報告する.
5. ユーザーがポータル上でベンチマーク結果を確認する.

## Links
- [ISUCON過去問環境をAWSで再現するための一式まとめ](https://github.com/matsuu/aws-isucon)
- [AWS - リージョンとゾーン](https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/UserGuide/using-regions-availability-zones.html)
- [AWS - VPC](https://docs.aws.amazon.com/ja_jp/vpc/latest/userguide/what-is-amazon-vpc.html)
- [AWS - インターネットゲートウェイ](https://docs.aws.amazon.com/ja_jp/vpc/latest/userguide/VPC_Internet_Gateway.html)
- [AWS - サブネット](https://docs.aws.amazon.com/ja_jp/vpc/latest/userguide/VPC_Subnets.html)
- [AWS - EC2](https://docs.aws.amazon.com/ja_jp/AWSEC2/latest/UserGuide/concepts.html)