# Product Requirements Document (PRD) - Automação da Compra de VR/VA

Este documento define as tarefas principais a serem implementadas para desenvolver o aplicativo de automação da compra de Vale Refeição (VR), conforme especificações dos arquivos `instrucoes.md`, `leiame.md` e `QWEN.md`.

> ⚠️ **Nota de Confidencialidade Importante**: Os nomes dos colaboradores não estão disponíveis nas planilhas por motivos de sigilo. Todas as referências aos colaboradores serão feitas exclusivamente através do campo **MATRICULA** como identificador único.

## Tarefas de Implementação

### 1. Fundação e Estrutura do Projeto

- [x] **Configurar dependências do projeto**: Adicionar bibliotecas necessárias para manipulação de arquivos Excel (ex: `github.com/xuri/excelize/v2`) ao `go.mod`.
- [x] **Definir estrutura de pacotes internos**: Criar os pacotes sugeridos (`internal/excel`, `internal/calculo`, `internal/modelo`, `internal/validacao`) para organizar a lógica da aplicação.
- [x] **Criar estruturas de dados (structs)**: Definir as structs `Colaborador` e `Periodo` para representar os dados dos funcionários e seus afastamentos/férias.
- [x] **Implementar sistema de tratamento de erros**: Criar um tipo de erro customizado (`ErroProcessamento`) para lidar com falhas durante o processamento das planilhas.

### 2. Leitura e Processamento de Dados

- [x] **Implementar função para ler planilhas Excel**: Criar uma função no backend (`LerPlanilha`) que utilize uma biblioteca de manipulação de Excel para abrir e ler os arquivos `.xlsx`.
- [~] **Implementar seleção de diretório no frontend**: Adicionar uma funcionalidade na interface React para que o usuário possa selecionar a pasta contendo as planilhas de entrada.
- [x] **Consolidar bases de dados**: Desenvolver a lógica para ler e combinar as informações de todas as 5 planilhas separadas (Ativos, Férias, Desligados, Base Cadastral, Base Sindicato) em uma única estrutura de dados em memória.
- [~] **Validar integridade dos dados**: Implementar validações para verificar:
    - Consistência de datas (admissão, desligamento, férias, afastamentos).
    - Preenchimento de campos obrigatórios.
    - Formatação correta dos dados (ex: formato de datas).
    - Relacionamento entre as planilhas através da matrícula.

> ⚠️ **Nota de Confidencialidade**: Os nomes dos colaboradores não estão disponíveis nas planilhas por motivos de sigilo. Para referenciar os colaboradores, utilizaremos exclusivamente o campo **MATRICULA** como identificador único.

### 3. Regras de Negócio e Cálculos

- [x] **Implementar regras de exclusão**: Desenvolver a lógica para filtrar e remover colaboradores não elegíveis da base consolidada, como:
    - Diretores, estagiários e aprendizes.
    - Profissionais afastados (licença maternidade, etc.).
    - Profissionais que atuam no exterior.
- [x] **Aplicar regras de datas "quebradas"**: Implementar o cálculo proporcional para colaboradores admitidos ou desligados no meio do mês.
- [x] **Implementar cálculo de dias úteis por sindicato**: Criar a lógica para calcular os dias úteis relevantes para cada colaborador, considerando:
    - O calendário de dias úteis de seu sindicato.
    - Férias (parciais ou integrais).
    - Outros afastamentos.
    - Data de desligamento.
- [x] **Implementar regra de desligamento**: Aplicar a regra específica para desligamentos: se comunicado até o dia 15, não considerar; se após o dia 15, calcular proporcionalmente.
- [x] **Calcular valores de VR**: Desenvolver a lógica para calcular o valor total de VR a ser concedido a cada colaborador, utilizando o valor definido para o seu sindicato.

> ⚠️ **Nota de Confidencialidade**: Como os nomes dos colaboradores não estão disponíveis nas planilhas, todas as referências aos colaboradores serão feitas exclusivamente através do campo **MATRICULA**.

### 4. Geração de Resultados e Interface

- [x] **Gerar planilha de resultado**: Criar uma função para gerar a planilha final no formato especificado (`./files/VR Mensal 05.2025.xls`), contendo:
    - Valor total de VR por colaborador (identificado pela **MATRICULA**).
    - Rateio 80%/20% (empresa/colaborador).
- [x] **Salvar planilha de resultado**: Implementar a funcionalidade para salvar a planilha gerada na pasta de Downloads do usuário.
- [x] **Implementar botão de análise no frontend**: Adicionar um botão "Fazer Análise" na interface que acione o processo de leitura, cálculo e geração da planilha.
- [X] **Exibir resultados na interface**: Criar uma área na interface para exibir o status do processamento e os resultados finais (opcional, mas desejável).
- [X] **Implementar interface de configuração**: Adicionar um menu de configuração para que o usuário possa informar chaves de API (OpenAI) ou configurar o acesso ao Ollama.

> ⚠️ **Nota de Confidencialidade**: A planilha de resultado utilizará a **MATRICULA** como identificador único do colaborador, uma vez que os nomes não estão disponíveis por motivos de sigilo.

### 5. Funcionalidades Adicionais

- [X] **Implementar funcionalidade de chat**: Criar uma interface de chat no frontend que permita ao usuário fazer perguntas sobre os dados processados, integrando com a API de IA.
- [X] **Adicionar validações da planilha**: Implementar as validações específicas indicadas na aba "validações" da planilha modelo.
- [X] **Considerar feriados**: Integrar um mecanismo (possivelmente através de uma API ou lista local) para considerar feriados estaduais e municipais no cálculo dos dias úteis.

### 6. Qualidade e Documentação

- [X] **Escrever testes unitários**: Criar testes para as funções críticas de cálculo e validação.
- [X] **Atualizar documentação**: Manter o `leiame.md` (ou `README.md`) atualizado com instruções de uso e desenvolvimento.
- [X] **Otimizar performance**: Garantir que o processamento seja eficiente, mesmo com grandes volumes de dados.
- [X] **Garantir segurança**: Proteger chaves de API e outros dados sensíveis.