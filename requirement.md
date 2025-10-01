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
8. Allow a web mode to launch a browser page that will enable to select database file if not already provided as command line parameter while launching the binary as a webserver
    - need REST endpoints to show
        - count of saved documents to the memory
        - count of each tool executed, like get_local_memory_document, get_all_documents, delete_local_memory_document, add_local_memory_document
        - view the existing documents
        - edit or add new document
        - mark or unmark a document as favourite
        - if content is edited or added, trigger the embedding and update / add embedding
        
    - UI will be developed using shadcn ui + react js + vite
    - it will have a layout, containing a sidebar
    - first page will be a dashboard page
        - showing, count of saved documents to the memory
        - count of each tool executed, like get_local_memory_document, get_all_documents, delete_local_memory_document, add_local_memory_document
    - second page will be
        - view the existing documents
        - edit or add new document
        - in the table view the user should be able to mark or unmark a document as favourite
        - if content is edited or added, trigger the embedding and update / add embedding



sample_docs = [
    "Albert Einstein proposed the theory of relativity, which transformed our understanding of time, space, and gravity.",
    "Marie Curie was a physicist and chemist who conducted pioneering research on radioactivity and won two Nobel Prizes.",
    "Isaac Newton formulated the laws of motion and universal gravitation, laying the foundation for classical mechanics.",
    "Charles Darwin introduced the theory of evolution by natural selection in his book 'On the Origin of Species'.",
    "Ada Lovelace is regarded as the first computer programmer for her work on Charles Babbage's early mechanical computer, the Analytical Engine."
]


https://docs.ragas.io/en/latest/getstarted/rag_eval/#up-next


sample_queries = [
    "Who introduced the theory of relativity?",
    "Who was the first computer programmer?",
    "What did Isaac Newton contribute to science?",
    "Who won two Nobel Prizes for research on radioactivity?",
    "What is the theory of evolution by natural selection?"
]

expected_responses = [
    "Albert Einstein proposed the theory of relativity, which transformed our understanding of time, space, and gravity.",
    "Ada Lovelace is regarded as the first computer programmer for her work on Charles Babbage's early mechanical computer, the Analytical Engine.",
    "Isaac Newton formulated the laws of motion and universal gravitation, laying the foundation for classical mechanics.",
    "Marie Curie was a physicist and chemist who conducted pioneering research on radioactivity and won two Nobel Prizes.",
    "Charles Darwin introduced the theory of evolution by natural selection in his book 'On the Origin of Species'."
]
