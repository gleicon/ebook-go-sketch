### Go e Estruturas de dados probabilisticas 

Este ebook é um resumo de como tenho utilizado estruturas de dados probabilisticas e Go.

Neste repositório tem o fonte do livro e os codigos que utilizei para testar e exemplificar ideias.

Gerei um preview do [Ebook aqui](ebook-go-sketches.pdf).

Se tiver alguma dúvida ou sugestão meu email é [gleicon@gmail.com](gleicon@gmail.com).

### Building

Para gerar um novo PDf a partir do arquivo `index.md` use o comando `make`.

### Artefatos para criar o PDF do EBook

O Makefile na raiz do repositório é usado para criar um PDF. Ele usa o Docker para executar o processo sem a necessidade de instalação de dependencias. Em casos comuns, executar o comando `make` vai criar um arquivo PDF novo.

O diretório buildtools tem o Dockerfile para criar a imagem que uso para converter de markdown para pdf com puppeteer, md-to-pdf (md2pdf) e o Google chrome. Esta imagem está em meu repositório no Docker Hub, então você só precisa recriar caso queira modificar o build.



Gleicon


