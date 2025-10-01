Local Memory Layer for Developer.
Maximum context for your coding agents, shared across IDEs, Projects and Teams.
Generated memories that scale with your codebase - programming concepts, bug fixes to business logic
Context relevant memory retrieval - bring back the right memory weather it is a bug fix, decision rule, or past discussion, so you never lose track of what matters.
Manage your memories - with add, search, list and delete memories tools.

This implementation provides a solid foundation for a lightweight, dependency-free vector database that can serve as an excellent MCP server for knowledge management tasks. The custom embedding algorithm, while simple, provides meaningful similarity matching for text documents without requiring external AI models.


The tool will be developed with golang that will

1. store embedding in a file based vector database like golang chromem DB "github.com/philippgille/chromem-go"
2. search for the embedding based on the query string and similarity threshold.
3. Documents will be short content like a passage, will have tags for easy lookup.
    - allow user to mark some of the documents as favourite which will rank higher in searches
    - allow to store key value pair of properties for each document which can be extended to store more such attributes like favourite document
    - add a column to hold create date
4. I want to achieve this without using any LLM or embedding model use.
    - use statistical methods to create content vector
    - use the same method to vectorize the query and search in the local vector database
5. Write in golang, so that the binary can be used without having to download any dependencies.
6. Use zerologger to print logs with filename and line number.
7. Enventually use this application as a MCP (Model Context Protocol) stdio server that will store content and retrive content based on user query
Reference for stdio server - https://modelcontextprotocol.io/llms-full.txt
Reference for golang MCP SDK - https://github.com/modelcontextprotocol/go-sdk
