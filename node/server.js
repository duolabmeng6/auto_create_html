const express = require('express');
const fs = require('fs');
const path = require('path');

const app = express();
const PORT = 3000;

app.engine('html', require('ejs').renderFile);
app.set('view engine', 'html');
app.set('views', path.join(__dirname, 'views'));

// 渲染EJS模板并保存到dist目录的通用函数
const renderAndSave = (page, res) => {
    res.render(page, {}, (err, html) => {
        if (err) {
            res.status(500).send('Error rendering template');
        } else {
            const outputFilePath = path.join(__dirname, 'dist', `${page}.html`);
            fs.writeFile(outputFilePath, html, 'utf8', (err) => {
                if (err) {
                    res.status(500).send('Error writing file');
                } else {
                    res.sendFile(outputFilePath);
                }
            });
        }
    });
};

app.get('/:page.html', (req, res) => {
    const page = req.params.page;
    renderAndSave(page, res);
});

app.use(express.static(path.join(__dirname, 'dist')));


app.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});
