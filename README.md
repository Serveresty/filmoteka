<h1>Filmoteka</h1>
<p>Решение тестового задания на позицию Go-разработчик(Стажёр)</p>
<hr />
<h2>Задание</h2>
<p><img src="https://github.com/Serveresty/filmoteka/assets/99687697/af47f23a-27f9-4066-8939-8ebb79a4bda9" /></p>
<hr />
<h2>Решение</h2>
<p>Выполнение требований:</p>
<ul>
  <li>Язык разработки Golang</li>
  <li>Для хранения данных использована реляционная СУБД(PostgreSQL)</li>
  <li>Предоставлена <a href="https://github.com/Serveresty/filmoteka/tree/main/docs">спецификация Swagger 2.0 (директория "./docs")</a></li>
  <li>Для реализации http сервера использована стандартная библиотека "net/http"</li>
  <li>Логирование реализовано через стандартную библиотеку "log" <a href="https://github.com/Serveresty/filmoteka/blob/main/pkg/logger/logger.go">(реализация logger)</a></li>
  <li>Написан Dockerfile <a href="https://github.com/Serveresty/filmoteka/tree/main/build">("./build")</a> и docker-compose.yaml <a href="https://github.com/Serveresty/filmoteka/tree/main/deployments">(./deployments)</a></li>
</ul>

<p>Разработанный функционал:</p>
<ul>
  <li>Все методы работы с фильмами</li>
  <li>Все методы работы с актёрами</li>
  <li>API закрыт авторизацией</li>
  <li>Имеется 2 роли: admin и user</li>
  <li>Все методы работы с пользователями</li>
</ul>

<hr />
<p>P.S. Было проведено ручное тестирование функционала (на написание unit-тестов нехватило времени и опыта, в частности использование fakeDatabase)</p>
