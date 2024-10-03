# Mirroring websites

When mirroring websites, we need to ensure we retrieve resources needed by the website, 
as defined by linked resources in the page. 
Websites, typically, link to other resources using HTML and CSS as defined in the following sections.

## CSS

### Linking to External Content in CSS

CSS primarily styles elements on a webpage, but it also offers limited ways to link to external resources. 

**1. `@import` Rule:**

*   **Purpose:** Imports external stylesheets into the current stylesheet.
*   **Syntax:** `@import "stylesheet_url";` or `@import url("stylesheet_url");`
*   **Notes:**
    *   Can import CSS files or other text files containing CSS rules.
    *   Can use relative or absolute URLs.

**2. `url()` Function:**

*   **Purpose:** Links to external resources within various CSS properties.
*   **Syntax:** `property: url("resource_url");`
*   **Notes:**
    *   Supports relative and absolute URLs.
    *   Can link to various file types (images, fonts, etc.).

## URL

A typical URL (Uniform Resource Locator) consists of several parts, each with a specific function. 
Here's an example of a URL with all possible parts:

```text
http://www.example.com:8080/path/to/resource?search=query&sort=asc&page=2#fragment
```

**Protocol**: The protocol specifies the communication method used to access the resource. 
In this case, it's HTTP (Hypertext Transfer Protocol). Other common protocols include HTTPS (HTTP Secure) and FTP (File Transfer Protocol).

**Domain**: The domain is the unique name of the website. 
It includes the subdomain (www), the second-level domain (example), and the top-level domain (.com).

**Subdomain**: A subdomain is an extension of a domain name, often used for organization or geographical purposes. 
In this example, "www" is the subdomain.

**Port**: The port is an optional number that specifies the communication endpoint within the host. 
The default port for HTTP is 80, but in this case, the URL is explicitly using port 8080.

**Path**: The path is a series of segments that identify the resource on the server. 
In this example, the path is "/path/to/resource".

**Query String**: The query string provides additional information about the resource and 
is appended to the path with a "?". In this example, the query string is "search=query&sort=asc&page=2". 
The "&" separates each parameter, and each parameter consists of a name (key) and value pair.

**Fragment**: The fragment, also known as the "anchor", 
is optional and specifies a specific location within the resource. 
In this example, the fragment is "#fragment".

## HTML elements linking to external content:

The following is a list of HTML elements commonly used to link to external content:

### 1. `<a>` (Anchor) - Hyperlinks
Purpose: Creates a clickable link to another page or resource.
Example:

```html
<a href="https://www.example.com">Visit our website</a>
```

### 2. `<img>` (Image) - Images
Purpose: Displays an image from an external source.
Example:

```html
<img src="https://www.example.com/image.jpg" alt="Example image">
```

### 3. `<video>` - Videos
Purpose: Embeds a video from an external source.
Example:

```html
<video src="https://www.example.com/video.mp4" controls></video>
```

### 4. `<audio>` - Audio
Purpose: Embeds an audio file from an external source.
Example:

```html
<audio src="https://www.example.com/audio.mp3" controls></audio>
```

### 5. `<iframe>` - Embedded Content
Purpose: Embeds content from another website or application.
Example:

```html
<iframe src="https://www.example.com/embedded-content" width="800" height="600"></iframe>
```

### 6. `<link>` (Stylesheet) - External Stylesheets
Purpose: Links an external stylesheet to the HTML document for styling.
Example:

```html
<link rel="stylesheet" href="https://www.example.com/styles.css">
```

### 7. `<script>` - External Scripts
Purpose: Links an external JavaScript file for functionality.
Example:

```html
<script src="https://www.example.com/script.js"></script>
```

### 8. `<object>` - External Content (Legacy)
Purpose: Includes external content like multimedia, applications, or other objects.
Example:

```html
<object data="https://www.example.com/object.swf" type="application/x-shockwave-flash"></object>
```
