# Segurança

Este documento descreve as práticas de segurança implementadas no BrxAgente-desafio4 para proteger chaves de API e outros dados sensíveis.

## Armazenamento de Chaves de API

As chaves de API são armazenadas com os seguintes princípios:

1. **Criptografia**: As chaves de API são criptografadas antes de serem salvas no disco.
2. **Armazenamento Seguro**: As chaves criptografadas são armazenadas em um diretório protegido no perfil do usuário (`.brxagente`).
3. **Chave de Criptografia**: Uma chave de criptografia AES-256 é gerada automaticamente e armazenada separadamente com permissões restritas.

## Implementação

O pacote `internal/security` implementa as seguintes funcionalidades:

- `encryption.go`: Funções para criptografar e descriptografar dados usando AES-256.
- `key_management.go`: Funções para gerar, salvar e carregar a chave de criptografia.
- `secure_string.go`: Tipo que representa uma string segura (criptografada em memória).

## Comunicação com APIs Externas

A comunicação com APIs externas segue estas práticas:

1. **HTTPS**: Todas as comunicações são feitas exclusivamente por HTTPS.
2. **Validação de Certificados**: Os certificados SSL são validados automaticamente.
3. **Tratamento de Erros**: Erros de comunicação são tratados adequadamente sem expor informações sensíveis.

## Revisões de Código

O código é regularmente revisado para identificar possíveis vazamentos de informações sensíveis, seguindo estas diretrizes:

1. Nenhuma chave de API é armazenada em texto plano no código.
2. Nenhuma chave de API é registrada em logs.
3. As chaves de API são tratadas como dados sensíveis em todo o código.