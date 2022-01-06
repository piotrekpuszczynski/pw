# pw
programowanie współbieżne

każda lista zawiera rozwiązanie zadania w języku go i ada.
* listy 1 nie wrzucałem - druga lista ją rozszeża

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

## lista 2

Program ma być uruchamiany z parametrami: n, d, k, b, h, klusownik

### a)
Rozszerz system zaimplementowany w zadaniu z poprzedniej listy w taki sposób, aby można w nim dodać b  krawędzi skierowanych postaci (i,j), gdzie i>j, oraz ustalić parametr h, oznaczający czas życia pakietu rozumiany jako największa liczba jego transferów od wierzchołka do wierzchołka. W grafie mogą występować cykle, więc jeśli pakiet w h krokach nie dotrze do celu, to  drukowany jest komunikat o jego śmierci i znika z systemu.
### b)
Dodaj wątek kłusownika, który co pewien czas budzi się, kontaktuje się z wątkiem losowo wybranego wierzchołka i umieszcza w nim pułapkę na jeden pakiet.  Jeśli pakiet dotrze do wierzchołka z zastawioną pułapką, to drukowany jest komunikat, że wpadł on w pułapkę i pakiet znika z systemu wraz z pułapką, w którą wpadł. 

## lista 3

Zakładamy, że mamy graf n wierzchołków, w którym krawędzie są nieskierowane. 
Krawędź między wierzchołkami i a j oznaczamy: {i,j}.
Listę sąsiadów wierzchołka i oznaczamy: N(i).
Podobnie jak w poprzednich zadaniach zakładamy, że w grafie istnieje ścieżka Hamiltona złożona z krawędzi postaci {v, v+1} (dzięki czemu graf jest spójny), oraz pewna liczba d dodatkowych krawędzi (skrótów). 
Należy zaimplementować wykonywanie protokołu routingu podobnego do znanego protokołu RIP, zgodnie z poniższymi wskazówkami.
* Każdy wierzchołek i zawiera zmienną reprezentującą tzw. routing table (oznaczaną przez Ri), która dla każdego wierzchołka j, różnego od i, zawiera następujące dane:
# * Ri[j].nexthop - wierzchołek ze zbioru N(i) (tj. sąsiad i) leżący na najkrótszej, znanej wierzchołkowi i, ścieżce p od i do j, oraz
# * Ri[j].cost - długość tej ścieżki p.
* Początkowo  każdy wierzchołek i zna swoich bezpośrednich sąsiadów N(i) i wie o istnieniu krawędzi postaci {v,v+1}. Zatem, 
dla jN(i),  początkowo Ri[j].cost=1  i  Ri[j].nexthop=j, a
dla jN(i), Ri[j].cost=|i-j|  oraz
Ri[j].nexthop=i+1, jeśli i<j, albo 
Ri[j].nexthop=i-1, jeśli j<i.
* Ponadto, dla każdegoRi[j], istnieje flaga Ri[j].changed (początkowo ustawiona na true).
* W każdym wierzchołku i działają dwa współbieżne wątki:
Senderi  oraz
Receiveri
* Oba te wątki mają współbieżny dostęp do routing table Ri. W Go można zaimplementować Ri jako stateful goroutine a w Adzie jako zmienną protected.
* Co pewien czas Senderi budzi się i jeśli istnieją jakieś j, gdzie Ri[j].changed=true, to tworzy pakiet z ofertą, do którego dodaje pary (j, Ri[j].cost) dla wszystkich takich j, ustawiając Ri[j].changedna false, a następnie wysyła ten pakiet do każdego swojego sąsiada z N(i).
* Wątek Receiveri oczekuje na pakiet z ofertą od jakiegoś sąsiada z N(i). Gdy taki pakiet otrzymuje od jakiegoś sąsiada l, to dla każdej pary (j, costj)z takiego pakietu:
wylicza newcosti,j=1+costj,
jeśli newcosti,j<Ri[j].cost to ustawia nowe wartości:
Ri[j].cost=newcost,
Ri[j].nexthop=l,
Ri[j].changed=true,
* Oba wątki drukują stosowne komunikaty o wysyłanych i otrzymywanych pakietach oraz zmianach w w routing table.
