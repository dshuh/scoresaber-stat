# template
template은 Golang REST API Template 프로젝트입니다.

## 선행작업
amigo 프로젝트를 $GOPATH/src에 다음과 같이 clone한다.

> $ cd $GOPATH/src

> $ git clone http://stove-gitlab.sginfra.net/amigo/amigoapp.git
> $ git clone http://stove-gitlab.sginfra.net/amigo/amigoaws.git
> $ git clone http://stove-gitlab.sginfra.net/amigo/amigocluster.git
> $ git clone http://stove-gitlab.sginfra.net/amigo/amigodb.git
> $ git clone http://stove-gitlab.sginfra.net/amigo/amigonet.git
> $ git clone http://stove-gitlab.sginfra.net/amigo/amigorpc.git
> $ git clone http://stove-gitlab.sginfra.net/amigo/amigoutil.git

template 프로젝트를 clone한다.

> $ git clone http://stove-gitlab.sginfra.net/backend/template.git

template 프로젝트로 이동해 의존성 패키지를 가져온다.

> $ cd template

> $ go get -v ./...

에러메세지가 출력될 경우 다시 한 번 시도한다.

## 테스트

테스트는 다음과 같이 수행한다.

> $ go test ./...

## 빌드하기

터미널을 열고 프로젝트 root 디렉토리에서 다음과 같이 실행한다.

- MAC

> $ make

빌드가 성공하면 다음과 같이 실행해 본다.

> $ cd bin

> $ ./template

- Windows

빌드 및 실행을 다음과 같이 실행한다.

> $ sh build.sh


빌드된 바이너리의 git revision과 빌드 날짜를 알고 싶다면 다음과 같이 확인할 수 있다.

> $ ./template -v
