# ResponseFormatter - Exemplos de Uso

## Visão Geral

O ResponseFormatter foi implementado para melhorar a formatação consistente de respostas no Chat Avançado, conforme especificado na issue #69.

## Tipos de Resposta Suportados

### 1. PolicyResponse - Consultas de Política
```
## 📋 Consulta de Política

**Pergunta:** Diretores têm direito a VR?

**Resposta:**
Não, diretores não têm direito ao Vale Refeição segundo a política VR_001.

**Fonte:** Manual de Recursos Humanos v2.1
**Confiança:** 95% (Muito Alta)

💡 *Esta resposta é baseada nas políticas oficiais da empresa.*
```

### 2. DataResponse - Dados Processados
```
## 📊 Consulta de Dados Processados

**Pergunta:** Quantos colaboradores foram processados?

**Resultado:**
Total de 150 colaboradores processados

**Estatísticas:**
- Processados: 150
- Pendentes: 5
- Erro: 2

*Dados baseados no último processamento: 12/09/2025 15:30*
```

### 3. CalculationResponse - Cálculos de VR
```
## 🧮 Cálculo de VR

**Cenário:** Calcular VR para colaborador admitido dia 20

**Resultado:** R$ 350,00

**Aplicação da Regra:**
Aplicar regra de data quebrada: proporcional aos dias úteis

**Cálculo:**
22 dias úteis - 15 dias trabalhados = R$ 450,00 * (15/22) = R$ 306,82

📖 **Política aplicada:** Política VR_005 - Datas Quebradas
```

### 4. ErrorResponse - Tratamento de Erros
```
## ❌ Erro no Processamento

**Erro:** Dados de colaborador não encontrados

**Sugestões:**
- Verificar se a matrícula está correta
- Tentar novamente em alguns minutos

💭 *Se o problema persistir, verifique os dados de entrada ou entre em contato com o suporte.*
```

### 5. WhatIfResponse - Cenários Hipotéticos
```
## 🤔 Análise Hipotética

**Cenário:** E se o colaborador fosse admitido dia 10?

**Resultado da Simulação:**
O valor seria R$ 450,00 (valor integral)

**Detalhes do Cálculo:**
Admissão dia 10: 22 dias úteis = valor integral

🔮 *Esta é uma simulação baseada nas políticas atuais. Resultados reais podem variar.*
```

## Como Usar

### No VRAgent
```go
// O ResponseFormatter é automaticamente integrado no VRAgent
agent := NewVRAgent(config, chatService)

// Métodos de formatação disponíveis:
response := agent.FormatPolicyResponse(question, answer, source, confidence)
response := agent.FormatDataResponse(question, answer, data, stats)
response := agent.FormatCalculationResponse(question, result, rule, calc, policy)
response := agent.FormatErrorResponse(errorMsg, suggestions)
response := agent.FormatWhatIfResponse(question, answer, calculation)
```

### Standalone
```go
formatter := NewResponseFormatter(nil) // Configuração padrão

data := ResponseData{
    Question: "Sua pergunta aqui",
    Answer: "Resposta aqui",
    // ... outros campos
}

formattedResponse := formatter.Format(PolicyResponse, data)
```

### Configuração Personalizada
```go
config := &FormatterConfig{
    UseEmojis: false,        // Desabilitar emojis
    DetailLevel: "minimal",  // Nível mínimo de detalhes
    IncludeFooter: false,    // Não incluir rodapé
    CompactMode: true,       // Modo compacto
}

formatter := NewResponseFormatter(config)
```

## Recursos Implementados

✅ Templates personalizáveis por tipo de resposta  
✅ Sistema de configuração flexível  
✅ Processamento de templates com funções helper  
✅ Formatação contextual (compacto, detalhado, técnico)  
✅ Integração automática com classificador de perguntas  
✅ Remoção/adição de emojis baseada na configuração  
✅ Níveis de confiança em texto (Muito Alta, Alta, etc.)  
✅ Timestamps automáticos  
✅ Suporte completo a testes unitários  

## Performance

O sistema foi otimizado para:
- Templates compilados uma única vez
- Cache de configurações
- Processamento eficiente de texto
- Mínimo overhead na formatação

## Extensibilidade

- Fácil adição de novos tipos de resposta
- Templates totalmente customizáveis
- Sistema de plugins para formatação específica
- Suporte a diferentes contextos de saída