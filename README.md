# keeper


Obs a serem feitas:

1- Dentro da pasta de infra tem todas as inicializações de infra do projeto, pq fazer assim? se precisar trocar um banco, não vai mudar nada

<img width="540" height="390" alt="image" src="https://github.com/user-attachments/assets/935951f9-c630-4e27-ad58-9a81589b8b51" />


2- Todo projeto é no padrão singleton, retornando a estrutura privada, de acordo com o DDD a camada do dominio não precisa conhecer sua implementação, ou seja, quando for implementar a interface, não preciso necessariamente implementar a interface Repository, posso ter duas interfaces diferentes desde que implemente o que for preciso

<img width="761" height="570" alt="image" src="https://github.com/user-attachments/assets/983b8ccb-e86e-4e71-868d-312f7e1cfcf7" />

3- não sei se é vantagem ou não, porém adiciono em cima dos arquivos um go:generate, facilita mt na hora de criar mock, pq se eu adicionar uma interface não preciso ficar procurando o que adicionei

<img width="919" height="335" alt="image" src="https://github.com/user-attachments/assets/e2455023-4eaf-46f8-a9b0-19162309d74d" />

