# Telegram bot

Telegram bot for saving bookmarks (notes, links, etc.)

## Project structure

- app - Main application
- clients - Different clients for services (API)
- cmd - Entry
- config - Configuration
- pool - Structure containing various services
- store - DAL
- units - Services

## Configuration
The configuration is specified using the config startup parameter, 
__config.yaml__ is used by default. So far, only the yaml driver is available.

