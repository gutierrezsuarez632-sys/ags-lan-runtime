# ags-lan-runtime

Motor de ejecución para **Aegis Framework**, un lenguaje orientado a dominio diseñado para generar estructuras de proyectos basadas en **DDD (Domain-Driven Design)** y **arquitectura hexagonal**.

Este repositorio contiene el **core del lenguaje y el runtime**, el cual será expuesto a través de un **daemon** consumido por herramientas externas como la CLI oficial.

---

## 🚀 Visión

AGS Language busca automatizar la creación de estructuras de proyectos complejos mediante un enfoque declarativo y orientado al dominio.

En lugar de escribir código repetitivo, defines la intención:

- Contextos delimitados (Bounded Contexts)
- Subcontextos
- Arquitectura (hexagonal, etc.)
- Configuración por lenguaje

Y el sistema genera la estructura automáticamente.

---

## 🧠 Filosofía

Este proyecto sigue estos principios:

- **Domain First**: el lenguaje representa conceptos de negocio, no detalles técnicos  
- **Declarativo > Imperativo**: defines *qué* quieres, no *cómo hacerlo*  
- **Extensible**: adaptable a múltiples lenguajes (Go, Node, PHP, etc.) Aunque su primer lanzamaiento estará disponible solo para Node (TypeScript) y PHP  
- **Desacoplado**: runtime independiente de CLI o interfaces externas  
---