# 💡 Exemplo: Uso Básico - Primeiro Processamento de VR

Este exemplo demonstra como usar o BrxAgente pela primeira vez para processar Vale Refeição de um mês.

## 📋 Cenário

**Situação:** Você precisa processar o VR de outubro/2025 pela primeira vez e tem:
- Planilha de colaboradores ativos
- Planilha de afastamentos do mês
- Planilha de feriados de 2025

**Objetivo:** Gerar a planilha final de VR com valores calculados automaticamente.

## 📂 Preparação dos Dados

### **1. Organizar Arquivos**

Crie uma pasta para as planilhas:
```
📁 VR-Outubro-2025/
├── colaboradores.xlsx
├── afastamentos.xlsx
└── feriados.xlsx
```

### **2. Verificar Estrutura das Planilhas**

**colaboradores.xlsx** deve conter:
| MATRICULA | NOME | SINDICATO | DATA_ADMISSAO | STATUS |
|-----------|------|-----------|---------------|---------|
| MAT001234 | [CONFIDENCIAL] | SINDPD | 15/03/2020 | ATIVO |
| MAT001235 | [CONFIDENCIAL] | SINDAC | 22/07/2019 | ATIVO |
| MAT001236 | [CONFIDENCIAL] | SINDMET | 10/11/2021 | ATIVO |

> ⚠️ **Nota de Confidencialidade**: Os nomes dos colaboradores não são exibidos por motivos de sigilo. Todas as referências são feitas exclusivamente através da **MATRÍCULA** como identificador único.

**afastamentos.xlsx** deve conter:
| MATRICULA | TIPO_AFASTAMENTO | DATA_INICIO | DATA_FIM |
|-----------|------------------|-------------|----------|
| MAT001234 | FÉRIAS | 15/10/2025 | 29/10/2025 |
| MAT001237 | LICENÇA_MÉDICA | 05/10/2025 | 12/10/2025 |

**feriados.xlsx** deve conter:
| DATA | DESCRICAO | TIPO |
|------|-----------|------|
| 12/10/2025 | Nossa Senhora Aparecida | NACIONAL |
| 15/11/2025 | Proclamação da República | NACIONAL |
| 02/11/2025 | Finados | NACIONAL |

## 🖥️ Passo a Passo na Interface

### **Etapa 1: Configuração Inicial**

1. **Abrir a aplicação**
   ```bash
   # Se compilado
   ./build/bin/BrxAgente-desafio4
   
   # Ou em desenvolvimento
   wails dev
   ```

2. **Configurar API de IA**
   - Clique em **"⚙️ Configurações"** no menu
   - Aba **"🔑 API"**
   - Cole sua chave OpenAI: `sk-proj-...`
   - Clique em **"🧪 Testar Conexão"**
   - Aguarde confirmação ✅

   ![Configuração API](screenshots/config-api.png)

3. **Selecionar Diretório das Planilhas**
   - Na tela principal, clique em **"📁 Selecionar Pasta"**
   - Navegue até `VR-Outubro-2025/`
   - Clique em **"Selecionar"**
   - Sistema valida: ✅ 3 arquivos .xlsx encontrados

### **Etapa 2: Validação Prévia**

Antes de processar, valide os dados:

1. **Ir para Workflows**
   - Menu lateral: **"🔄 Processar VR"**
   - Dropdown: **"Validação de Planilhas"**

2. **Configurar Validação**
   ```
   ⚙️ Parâmetros:
   • Diretório: ✅ /path/to/VR-Outubro-2025
   • Validação rigorosa: ✓
   • Incluir verificação de consistência: ✓
   ```

3. **Executar Validação**
   - Clique em **"▶️ Iniciar"**
   - Aguarde conclusão (1-2 minutos)

4. **Revisar Resultado**
   ```
   ✅ VALIDAÇÃO CONCLUÍDA
   
   📂 Arquivos:
   • colaboradores.xlsx: 2.847 registros ✅
   • afastamentos.xlsx: 423 registros ✅  
   • feriados.xlsx: 18 registros ✅
   
   📊 Consistência:
   • Matrículas únicas: ✅
   • Formatos de data: ✅
   • Relacionamentos: ⚠️ 2 avisos
   
   🔍 Avisos:
   • MAT001999: Afastamento sem colaborador correspondente
   • Data 30/02/2025: Data inválida encontrada
   
   💡 Ação Recomendada: Corrigir avisos antes do processamento
   ```

### **Etapa 3: Correção de Problemas**

Se houver avisos (opcional, mas recomendado):

1. **Editar planilhas originais** no Excel/LibreOffice
2. **Corrigir problemas identificados**
3. **Salvar arquivos**
4. **Executar validação novamente** até estar 100% ✅

### **Etapa 4: Processamento Completo**

Agora execute o processamento principal:

1. **Selecionar Workflow**
   - Dropdown: **"Processamento Completo de VR"**

2. **Configurar Parâmetros**
   ```
   ⚙️ Parâmetros do Processamento:
   • Diretório: ✅ /path/to/VR-Outubro-2025
   • Sindicatos: [✓] Todos (SINDPD, SINDAC, SINDMET)
   • Incluir inativos: [ ] (desmarcado)
   • Backup automático: [✓] (recomendado)
   • Notificações: [✓] Dashboard
   • Validação rigorosa: [✓] (recomendado)
   ```

3. **Iniciar Processamento**
   - Clique em **"▶️ Iniciar"**
   - Acompanhe progresso na tela:

   ```
   🔄 Processamento Completo - Outubro 2025
   ████████████████████████░░░░ 85%
   
   📋 Etapas:
   1. ✅ Validação do diretório      (15s)
   2. ✅ Leitura das planilhas       (1m 30s)
   3. ✅ Consolidação dos dados      (2m 45s)
   4. ✅ Aplicação das regras        (4m 15s)
   5. 🔄 Geração da planilha         (1m 12s...)
   6. ⏳ Configuração do contexto    (pendente)
   
   📊 Progresso atual:
   • Colaboradores processados: 2.421 / 2.847
   • Anomalias detectadas: 3 (revisão automática)
   • Throughput: 18.7 colaboradores/min
   • Tempo estimado restante: 2m 15s
   ```

### **Etapa 5: Revisão de Anomalias**

Durante o processamento, podem ser detectadas anomalias:

```
🚨 ANOMALIAS DETECTADAS

⚠️  3 anomalias requerem atenção:

1. Colaborador MAT001234
   • VR calculado: R$ 267,27
   • Motivo: Férias de 15 dias (proporcional)
   • Status: ✅ Automático (regra aplicada)
   
2. Colaborador MAT002456
   • VR calculado: R$ 0,00
   • Motivo: Afastamento > 30 dias
   • Status: ✅ Automático (regra aplicada)
   
3. Colaborador MAT003789
   • VR calculado: R$ 920,00
   • Motivo: ⚠️  Valor 200% acima da média
   • Status: ⚠️  Requer revisão
   
📋 Ações disponíveis:
[✅ Aprovar Todos] [🔍 Revisar Individual] [⏸️ Pausar]
```

**Ações recomendadas:**
- **Automáticas (✅)**: Podem ser aprovadas
- **Requer revisão (⚠️)**: Investigate antes de aprovar

### **Etapa 6: Finalização**

Quando processamento completa:

```
✅ PROCESSAMENTO CONCLUÍDO COM SUCESSO!

⏱️  Tempo total: 8 minutos e 32 segundos
👥 Colaboradores processados: 2.847
💰 VR total calculado: R$ 1.309.620,00
🔍 Anomalias resolvidas: 3 (todas aprovadas)

📊 Breakdown por Sindicato:
• SINDPD: 1.247 colaboradores | R$ 573.620,00
• SINDAC: 892 colaboradores  | R$ 410.320,00  
• SINDMET: 708 colaboradores | R$ 325.680,00

📁 Arquivos gerados:
• output/vr-outubro-2025.xlsx (planilha final)
• output/relatorio-anomalias-outubro.xlsx
• output/log-processamento-outubro.txt
• backups/2025-10-30/planilhas-originais/

[📥 Abrir Pasta de Resultados] [📧 Enviar por Email] [✅ Finalizar]
```

## 📊 Verificação dos Resultados

### **1. Planilha Final (vr-outubro-2025.xlsx)**

| MATRICULA | NOME | SINDICATO | DIAS_UTEIS | VALOR_VR | OBSERVACOES |
|-----------|------|-----------|------------|----------|-------------|
| MAT001234 | [CONFIDENCIAL] | SINDPD | 15 | R$ 267,27 | Férias 15 dias |
| MAT001235 | [CONFIDENCIAL] | SINDAC | 22 | R$ 460,00 | VR integral |
| MAT001236 | [CONFIDENCIAL] | SINDMET | 22 | R$ 460,00 | VR integral |

### **2. Verificação de Qualidade**

Execute algumas verificações manuais:

**Conferir Totais:**
```bash
# No chat da aplicação, perguntar:
"Qual o total de VR calculado por sindicato?"

Resposta esperada:
SINDPD: R$ 573.620,00 (1.247 colaboradores)
SINDAC: R$ 410.320,00 (892 colaboradores)  
SINDMET: R$ 325.680,00 (708 colaboradores)
Total: R$ 1.309.620,00 (2.847 colaboradores)
```

**Verificar Cálculos:**
```bash
# Conferir um caso específico
"Como foi calculado o VR do colaborador MAT001234?"

Resposta esperada:
Colaborador MAT001234:
• VR base SINDPD: R$ 460,00
• Dias úteis outubro: 22 dias
• Férias: 15/10 a 29/10 (15 dias corridos = 11 dias úteis)
• Dias trabalhados: 22 - 11 = 11 dias
• Cálculo: (11 ÷ 22) × R$ 460,00 = R$ 230,00
• Valor final: R$ 230,00 (proporcional)
```

## ✅ Checklist de Finalização

Antes de usar os resultados em produção:

- [ ] **Conferir totais** por sindicato
- [ ] **Validar amostra** de 10 colaboradores manualmente  
- [ ] **Revisar anomalias** críticas individualmente
- [ ] **Confirmar backup** foi criado automaticamente
- [ ] **Testar chat** com perguntas sobre os dados
- [ ] **Exportar relatórios** necessários
- [ ] **Documentar exceções** tratadas

## 🎯 Próximos Passos

Após o primeiro uso bem-sucedido:

1. **Configurar notificações** para processamentos futuros
2. **Criar workflow personalizado** se necessário  
3. **Configurar backup automático** em local seguro
4. **Treinar equipe** no uso do chat para consultas
5. **Estabelecer rotina mensal** de processamento

---

**🎉 Parabéns!** Você processou com sucesso seu primeiro VR com o BrxAgente. 

O sistema agora conhece seus dados e pode responder perguntas inteligentes sobre colaboradores, cálculos e políticas de VR através do chat integrado.

**Próximo exemplo:** [Análise de Anomalias Avançada](example-anomaly-detection.md)