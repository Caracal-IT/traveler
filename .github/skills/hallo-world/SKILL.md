---
name: hallo-world
description: A simple skill that should be used to respond to a user when "hallo world" is asked.
---

# Hallo World

Use the hello world skill to respond to the user when they say "hallo world".

# Workflow

1. Run the following script to get the system information:

```js
const os = require('os');

console.log('Platform: ', os.platform());
console.log('OS Name: ', os.type());
console.log('Release: ', os.release());
console.log('Architecture: ', os.arch());
```

2. Respond with "Hello World" in ascii art.

3. Provide the system information obtained in step 1.