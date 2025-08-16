# LEIAME

## Sobre

Este aplicativo permite analizar um conjunto de planilhas Excel com o objetivo de automatizar o processo mensal de compra de VR (Vale Refeição), garantindo que cada colaborador receba o valor correto, considerando ausências, férias e datas de admissão ou desligamento e calendário de feriados.

ver mais detalhes no arquivo "./files/instructions.md" 

Esta análise é feita usando recursos de inteligencia artificial, usando LLM.
O aplicativo permite configurar chaves de API para ferramentas online como OPEN AI (CHAT GPT), bem como permite também a analise a partir de modelos em execução local, usando o Ollama.

## Fonte de Dados
Os dados para esta analise estão no formato *.XLSX. Para iniciar a análise, o caminho para a pasta em que estes arquivos estão deve ser informado na interface gráfica.
As planilhas XSLX contem os dados dos funcionários. O código tenta encontrar uma relação entre as planilhas (por exemplo: Matrícula). Caso nao consiga, deverá retornar uma mensagem ao usuário explicando. 

## Tecnologias Utilizadas
O aplicativo foi desenvolvido usando o framework Wails: Golang no backend e React no frontend.
O framework utilizado para a parte de análise dos dados via Modelo de Linguagem é o LangChainGo, visto que a análise das planilhas deve ser executada no backend.

Para o Front End foi utilizado o React JS.

## Arquitetura do Projeto

O projeto segue uma arquitetura modular com os seguintes pacotes internos:

- `internal/excel`: Responsável pela manipulação de planilhas Excel
- `internal/calculo`: Contém a lógica de cálculo do VR
- `internal/modelo`: Define as estruturas de dados utilizadas no projeto
- `internal/validacao`: Implementa as regras de validação dos dados

## Características

* Permite a análise de um diretório com um conjunto de planilhas com os dados de funcionarios visando automatizar o processo mensal de compra de Vale Refeição.
* Permite informar a Chave de API, através de um menu de configuração, em que o usuário define se gostaria de utlizizar o OpenAI com ChatGPT, ou se pretende fazer esta analise localmente com um modelo configurado em um servidor Ollama rodando localmente.
* A analise é executada ao clicar num botão chamado "Fazer Análise". O resultado deve ser uma planilha de excel, conforme orientado no arquivo "./instrucoes.md". O aplicativo deve gerar este excel e deve gravar este arquivo na pasta de Downloads do usuário.
* Adicionalmente, apresenta uma interface de Chat em que o usuário pode fazer perguntas sobre os dados.
* A interface do aplicativo utiliza o idioma Português Brasileiro.

## Desenvolvimento

Para executar em modo de desenvolvimento, execute `wails dev` no diretório do projeto. Isso executará um servidor de desenvolvimento Vite
que fornecerá um hot reload muito rápido das suas alterações no frontend. Se você quiser desenvolver em um navegador
e tiver acesso aos seus métodos Go, também há um servidor de desenvolvimento que roda em http://localhost:34115. Conecte-se
a ele no seu navegador e você poderá chamar seu código Go a partir do devtools.

## Construir Executável

Para construir um pacote redistribuível em modo de produção, use `wails build`.
