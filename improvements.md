The current extractor implementation only retrieves the first page of results (e.g., the first 100 users or projects) from the Asana API due to hard limits imposed by the service. This is a functional correctness issue that prevents the solution from meeting the requirement of handling "thousands" of projects and employees.

Implementation Plan
To ensure the extractor retrieves a complete data set, we must implement a pagination loop based on the Asana API's cursor model:

Loop Structure: The API call will be wrapped in a loop that continues until all records are processed.

Cursor Check: After a successful API request, the program must inspect the response for the next_page object.

Next Request: If the next_page object contains an offset token, this token must be appended as a query parameter to the URL for the next API request.

Termination: The loop terminates when the JSON response returns a next_page value of null, signifying that the final page of data has been extracted.

This ensures the extractor is correct and capable of handling data volumes at enterprise scale.