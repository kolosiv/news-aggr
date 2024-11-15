document.getElementById('news-form').addEventListener('submit', function(event) {
    event.preventDefault();
    
    const startDate = document.getElementById('start-date').value;
    const endDate = document.getElementById('end-date').value;

    if (!startDate || !endDate) {
        alert("Both start and end dates are required.");
        return;
    }
    
    // Формируем URL для запроса с параметрами start и end
    const url = `/news/interval?start=${encodeURIComponent(startDate)}&end=${encodeURIComponent(endDate)}`;
    
    // Выполняем запрос к API
    fetch(url)
        .then(response => {
            if (!response.ok) {
                throw new Error('Error fetching news');
            }
            return response.json();
        })
        .then(data => {
            const newsList = document.getElementById('news-list');
            newsList.innerHTML = '';  // Очищаем список перед добавлением новых данных
            data.forEach(news => {
                const listItem = document.createElement('li');
                
                // Получаем дату из поля 'PubDate' и парсим ее
                const createdAtDate = new Date(news.PubDate);
                
                // Проверяем, чтобы дата была валидной
                if (isNaN(createdAtDate.getTime())) {
                    listItem.innerHTML = `Invalid Date for: ${news.Title}`;
                } else {
                    // Форматируем дату на dd-mm-yyyy и время на 24-часовой формат
                    const formattedDate = createdAtDate.toLocaleDateString('en-GB');  // формат dd-mm-yyyy
                    const formattedTime = createdAtDate.toLocaleTimeString('en-GB', { hour12: false });  // 24-часовой формат времени
                    
                    // Получаем название источника, если оно есть
                    const sourceName = news.SourceName ? `Источник: ${news.SourceName}` : 'Источник неизвестен';

                    listItem.innerHTML = `
                        <a href="${news.Link}" target="_blank">${news.Title}</a>
                        <p class="description">${news.Description}</p>
                        <p class="news-date">${formattedDate} ${formattedTime} - ${sourceName}</p>
                    `;
                }

                newsList.appendChild(listItem);
            });
        })
        .catch(error => {
            console.error('Error fetching news:', error);
            alert('There was an error fetching the news.');
        });
});
