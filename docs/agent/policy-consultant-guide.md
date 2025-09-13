# 📋 Guia do Consultor de Políticas

## Visão Geral
O **PolicyConsultantTool** é um sistema inteligente que fornece respostas precisas sobre políticas de Vale Refeição (VR) com base em uma base de conhecimento estruturada. Ele funciona como um especialista virtual em regulamentações de VR, capaz de responder perguntas complexas com citação de fontes oficiais.

## 📑 Índice

1. [Base de Conhecimento](#-base-de-conhecimento)
2. [Tipos de Consulta](#-tipos-de-consulta)
3. [Sistema de Classificação](#-sistema-de-classificação)
4. [Exemplos Práticos](#-exemplos-práticos)
5. [Níveis de Confiança](#-níveis-de-confiança)
6. [Limitações](#-limitações)

## 📚 Base de Conhecimento

### Políticas Disponíveis

O sistema possui conhecimento estruturado sobre as seguintes políticas:

#### **VR_001: Elegibilidade para Vale Refeição**
- Critérios de elegibilidade por categoria de colaborador
- Exclusões por tipo de contrato
- Exceções especiais aprovadas

#### **VR_002: Valor Base do Vale Refeição**
- Valores por sindicato
- Reajustes anuais
- Critérios de variação

#### **VR_003: Exclusões de Vale Refeição**
- Diretores estatutários
- Estagiários e aprendizes
- Colaboradores com jornada reduzida
- Casos especiais

#### **VR_004: Cálculo Proporcional**
- Regras para admissões no meio do mês
- Desligamentos e comunicados
- Licenças médicas e férias
- Afastamentos temporários

#### **VR_005: Calendário e Dias Úteis**
- Feriados nacionais, estaduais e municipais
- Diferenças por sindicato
- Recálculos para meses atípicos

#### **VR_006: Processamento e Prazos**
- Cronograma mensal
- Prazos de comunicação
- Exceções para processamento

## 🎯 Tipos de Consulta

### 1. **Simple (Consultas Diretas)**
Perguntas objetivas com resposta direta baseada na base de conhecimento.

**Exemplos:**
- "Diretores têm direito a VR?"
- "Qual o valor do VR para SINDPD?"
- "Estagiários recebem Vale Refeição?"

### 2. **Complex (Cenários com Múltiplos Fatores)**
Situações que envolvem múltiplas variáveis e regras combinadas.

**Exemplos:**
- "Colaborador do SINDAC, admitido dia 20, em licença médica de 10 dias, qual o VR?"
- "Desligamento comunicado dia 18 do mês, mas com férias programadas, tem direito?"

### 3. **WhatIf (Análises Hipotéticas)**
Simulações de cenários para planejamento ou entendimento de regras.

**Exemplos:**
- "E se um diretor tiver salário abaixo de R$ 15.000?"
- "O que acontece se comunicarmos desligamento depois do dia 15?"

### 4. **Compliance (Verificação Regulatória)**
Consultas sobre conformidade e aderência às normas.

**Exemplos:**
- "Estamos seguindo corretamente a regra de dias úteis?"
- "A exclusão de estagiários está conforme a legislação?"

## 🔍 Sistema de Classificação

O sistema classifica automaticamente as perguntas usando as seguintes categorias:

### **Palavras-chave de Identificação:**

#### **Políticas:**
- "direito", "pode", "elegível", "política"
- "permitido", "autorizado", "regra"
- "diretor", "estagiário", "aprendiz"

#### **Cálculos:**
- "calcular", "valor", "quanto", "como"
- "proporcional", "dias", "licença"
- "admissão", "desligamento", "férias"

#### **Dados:**
- "quantos", "total", "soma", "média"
- "colaboradores", "processados", "distribuição"
- "sindicato", "empresa", "setor"

#### **Cenários Complexos:**
- "se", "caso", "quando", "situação"
- "mas", "porém", "exceto", "além"
- "combinado", "junto", "simultâneo"

## 💡 Exemplos Práticos

### **Exemplo 1: Consulta Simple**
```
👤 Pergunta: "Diretores têm direito a VR?"

🤖 Resposta:
## 📋 Consulta de Política
**Pergunta:** Diretores têm direito a VR?

**Resposta:** Não, diretores não têm direito ao Vale Refeição conforme estabelecido na política da empresa.

**Detalhes:**
• Diretores Estatutários: Sempre excluídos
• Diretores Executivos: Sempre excluídos  
• Diretores de Área: Sempre excluídos

**Exceção:** Diretores com salário inferior a R$ 15.000 podem ser incluídos mediante aprovação especial da diretoria.

**Fonte:** Política VR_003 - Manual de RH v2.1, Seção 2.3
**Confiança:** 95% (Muito Alta)
```

### **Exemplo 2: Consulta Complex**
```
👤 Pergunta: "Colaborador do SINDAC admitido dia 25, como calcular o VR?"

🤖 Resposta:
## 🧮 Cálculo de VR - Cenário Complexo
**Cenário:** Admissão no meio do mês (SINDAC)

**Cálculo Proporcional:**
Para admissões após o dia 15, aplicar fórmula proporcional:

**Fórmula:** (Dias úteis trabalhados ÷ Total dias úteis) × VR base

**Exemplo Prático (Outubro 2025):**
• Total dias úteis SINDAC: 22 dias
• Dias trabalhados (25 a 31): 5 dias úteis
• VR base SINDAC: R$ 460,00
• **Cálculo:** (5 ÷ 22) × R$ 460,00 = R$ 104,55

**Observações:**
• Considerar apenas dias úteis efetivamente trabalhados
• Não incluir sábados, domingos ou feriados
• Verificar se há afastamentos no período

**Fonte:** Política VR_004 - Cálculo Proporcional v1.3
**Confiança:** 90% (Alta)
```

### **Exemplo 3: Consulta WhatIf**
```
👤 Pergunta: "E se comunicarmos um desligamento depois do dia 15?"

🤖 Resposta:
## 🎯 Análise Hipotética - Regra de Desligamento
**Cenário:** Comunicação de desligamento após dia 15

**Consequências:**
• **Direito integral ao VR:** Colaborador mantém direito ao valor completo do mês
• **Processamento normal:** Incluir no cálculo mensal normalmente
• **Impacto financeiro:** Empresa arca com valor total mesmo com trabalho parcial

**Regra Oficial:**
"Desligamentos comunicados após o dia 15 do mês garantem direito integral ao Vale Refeição do mês corrente"

**Recomendação:**
• Comunicar desligamentos preferencialmente até dia 15
• Para otimização de custos, planejar desligamentos para início do mês seguinte
• Exceções devem ser aprovadas pela gestão

**Fonte:** Política VR_006 - Processamento e Prazos v2.0
**Confiança:** 88% (Alta)
```

### **Exemplo 4: Consulta Compliance**
```
👤 Pergunta: "Estamos seguindo corretamente a exclusão de estagiários?"

🤖 Resposta:
## ✅ Verificação de Conformidade - Exclusão de Estagiários
**Status:** ✅ Conforme

**Verificação Realizada:**
• **Base Legal:** Lei nº 11.788/2008 (Lei do Estágio)
• **Política Interna:** VR_003 - Exclusões
• **Implementação Atual:** ✅ Correta

**Detalhes da Conformidade:**
• ✅ Estagiários não recebem VR (conforme Lei do Estágio)
• ✅ Aprendizes também excluídos (Lei 10.097/2000)
• ✅ Separação clara nos sistemas internos
• ✅ Documentação adequada nos contratos

**Recomendações:**
• Manter documentação atualizada
• Revisar novos contratos de estágio mensalmente
• Verificar classificação correta no sistema

**Fonte:** Política VR_003 + Legislação Federal
**Confiança:** 98% (Muito Alta)
```

## 🎯 Níveis de Confiança

O sistema indica o nível de confiança de cada resposta:

### **🔴 Baixa (0-50%)**
- Informações incompletas na base de conhecimento
- Cenários não cobertos pelas políticas atuais
- Necessidade de verificação manual

### **🟡 Média (51-75%)**
- Informações parciais disponíveis
- Cenários com múltiplas interpretações possíveis
- Recomendação de validação adicional

### **🟢 Alta (76-90%)**
- Informações claras e bem documentadas
- Cenários cobertos pelas políticas existentes
- Confiável para uso direto

### **🟢 Muito Alta (91-100%)**
- Informações totalmente documentadas
- Casos explicitamente cobertos na base de conhecimento
- Fonte oficial disponível e citada

## ⚠️ Limitações

### **O que o Sistema NÃO pode fazer:**
- ❌ Alterar políticas ou criar exceções
- ❌ Acessar dados de colaboradores específicos (por sigilo)
- ❌ Processar cálculos em tempo real sem contexto
- ❌ Interpretar leis que não estão na base de conhecimento
- ❌ Dar respostas sobre casos não documentados

### **O que Requer Validação Humana:**
- 📋 Casos excepcionais não previstos nas políticas
- 📋 Mudanças na legislação recentes
- 📋 Interpretações de contratos específicos
- 📋 Decisões que envolvem aprovações especiais

### **Quando Consultar o RH:**
- 🔍 Confiança abaixo de 75%
- 🔍 Casos que fogem do padrão
- 🔍 Necessidade de exceções ou alterações
- 🔍 Situações envolvendo questões legais complexas

## 🔄 Atualização da Base de Conhecimento

A base de conhecimento deve ser atualizada regularmente para manter a precisão:

### **Frequência Recomendada:**
- **Mensal:** Revisão de políticas ativas
- **Trimestral:** Validação de casos novos
- **Anual:** Atualização completa da base
- **Ad-hoc:** Quando há mudanças na legislação

### **Processo de Atualização:**
1. Identificação de novos casos ou mudanças
2. Documentação oficial das políticas
3. Teste com casos conhecidos
4. Validação das respostas
5. Deploy da versão atualizada

---

**O PolicyConsultantTool é uma ferramenta poderosa para consultas sobre políticas de VR, mas sempre deve ser usado em conjunto com validação humana para casos complexos ou críticos! 🚀**