# Windows Setup Guide

Hey devs! 👋

This is a **weekend side project** I made — a tiny Go tool that turns **plain English into SQL queries**. It works wonderfully with **MySQL Workbench** (and may work with phpMyAdmin, though I haven’t fully tested that).

Remember: this is **just for fun**, not a polished product.

## Prerequisites

* Make sure you have **[Go](https://golang.org/dl/) installed** on your system.
* Clone this repo to your local machine:

```bash
git clone https://github.com/cttricks/wb-agent.git
cd wb-agent
```

## Prepare Necessary Files

Before running the program, you need to copy a few important files to the **root of the project**:

1. [**Instruction.yaml**](/assets/Instruction.yaml)

   * Located in `assets/Instruction.yaml`.
   * This file is used to generate `Instruction.md` and contains **demo tables and settings**.
   * **Important:** If you make any changes to `Instruction.yaml`, delete the existing `Instruction.md` so it can regenerate correctly.

2. [**agent.ico**](/assets/agent.ico)

   * Located in `assets/agent.ico`.
   * Keeps the Windows notifications looking nice when the tool runs.

3. **.env file**

   * Create a `.env` in the root with the following keys:

```env
OPENAI_API_KEY=sk-proj-XXXXXXX
DATABASE_NAME=mydatabase
```

> Pro Tip: If this feels confusing, you can also download the **latest release ZIP**, extract it anywhere, and it already contains all these necessary files (except the Go source code).



## Customize Instruction.yaml

The AI relies on `Instruction.yaml` to understand your database. Think of it as **teaching the AI your DB structure**, because otherwise, garbage in → garbage out.

* Add your tables under `tables:` like this:

```yaml
tables:
  - name: table_name
    description: Describe briefly why this table exists
    columns:
      - name: column_name
        type: INT|VARCHAR|DATETIME
        notes: Primary key, auto-incremented, etc.
      ...
    usage:
      - "Use this table for <purpose> queries."
      - "Filter active users with `active = 1`."
      ...
    examples:
      - question: "Count how many users have a phone number."
        sql: |
          SELECT COUNT(*) FROM t_user WHERE phone IS NOT NULL;
      ...
```

* You can add as many tables as needed.
* **Important:** AI works on simple principles, so this structure is necessary for it to generate meaningful queries.



## Run the Program

Once you have all files ready:

```bash
go run .
```

This will initialize the program and start listening for your English-to-SQL commands.

## Open MySQL Workbench

1. Make sure the **query tab is blank**.
2. Type a plain English question, e.g.:

```
How many members are there in team
```

3. Press **ALT + A**

The AI will take your input and generate SQL, e.g.:

```sql
SELECT COUNT(*) as members FROM t_users WHERE active = 1;
```

* Execute the query in **Workbench** or try in **phpMyAdmin** (works best in Workbench).

## Notes & Contributions

* ⚠️ This is a **side project**. Bugs will happen!
* Found a bug and have a fix? **Raise a PR**.
* Just want to report it? I’ll try to push updates when I have free time.
* If you enjoy tinkering, a **STAR ⭐** on the repo is much appreciated!



Made with ☕, Go, and a weekend curiosity.