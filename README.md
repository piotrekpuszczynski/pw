# pw
programowanie współbieżne

każda lista zawiera rozwiązanie zadania w języku go i ada.

## lista 1

Program ma być uruchamiany z parametrami: n, d, b

### Program ma działać następująco:
* Generowany jest graf G dla podanych parametrów n i d, gdzie d skrótów generowane jest w sposób losowy.
* Graf G drukowany jest na terminalu tak aby przedstawić istniejące połączenia. (Zastanowić się nad tym jaki sposób prezentacji będzie najbardziej czytelny.)
* Uruchamiana jest symulacja systemu przesyłania pakietów po grafie G.
### System  przesyłania pakietów działa według następujących zasad:
* W jednym wierzchołku może przebywać tylko jeden pakiet.
* Co pewien losowy czas nadawca umieszcza w źródle (o ile jest ono puste) nowy pakiet indeksowany kolejną liczbą naturalną.
* Co pewien losowy czas odbiorca odbiera z ujścia pakiet (o ile jest co odebrać).
* Pakiet w wierzchołku i, po odczekaniu losowego czasu, wybiera losowo jeden wierzchołek j ze zbioru N(i) i czeka aż będzie mógł się do niego przemieścić. 
* Gdy pakiet p dotrze do wierzchołka i drukowany jest komunikat:
"pakiet p jest w wierzchołku i"
i jednocześnie p dodaje i do swojej listy odwiedzonych wierzchołków oraz i dodaje p do swojej listy obsłużonych pakietów.
* Gdy odbiorca odbierze pakiet p, drukuje komunikat:
"pakiet p został odebrany".     
* Po nadaniu kpakietów, nadawca kończy nadawanie.
* Gdy odbiorca odbierze ostatni (tj. k-ty)  pakiet, system kończy działanie i rozpoczyna się drukowanie raportów końcowych.
* W raportach końcowych pojawią się dwa wykazy:
dla każdego wierzchołka, lista kolejno obsłużonych przez niego pakietów, 
dla każdego pakietu, lista odwiedzonych przez niego wierzchołków  (ścieżka od źródła do ujścia).

## lista2

Program ma być uruchamiany z parametrami: n, d, b, k, h

### a)
Rozszerz system zaimplementowany w zadaniu z poprzedniej listy w taki sposób, aby można w nim dodać b  krawędzi skierowanych postaci (i,j), gdzie i>j, oraz ustalić parametr h, oznaczający czas życia pakietu rozumiany jako największa liczba jego transferów od wierzchołka do wierzchołka. W grafie mogą występować cykle, więc jeśli pakiet w h krokach nie dotrze do celu, to  drukowany jest komunikat o jego śmierci i znika z systemu.
### b)
Dodaj wątek kłusownika, który co pewien czas budzi się, kontaktuje się z wątkiem losowo wybranego wierzchołka i umieszcza w nim pułapkę na jeden pakiet.  Jeśli pakiet dotrze do wierzchołka z zastawioną pułapką, to drukowany jest komunikat, że wpadł on w pułapkę i pakiet znika z systemu wraz z pułapką, w którą wpadł. 

